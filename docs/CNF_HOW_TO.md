How to make CNF compatible with StoneWork
=========================================

This guide explains how to develop CNF using ligato framework and make
it interoperable and dynamically loadable by StoneWork.
For more information about how CNFs and StoneWork interact with each other
see [StoneWork Architecture][architecture].


1. Firstly, for CNF to be able to talk to StoneWork it has to load and init two plugins from the
   StoneWork repository: [CNF Registry][cnf-registry-plugin] and [Punt Manager][punt-manager-plugin].
   Since these plugins have to be imported and because StoneWork is inside a private gerrit repository,
   it is necessary to add StoneWork as a git submodule into the CNF repository:
   ```
   cnf-repo$ git submodule add https://github.com/PANTHEONtech/StoneWork submodule/stonework
   ```
   Then add into go.mod replace directive:
   ```
   replace go.pantheon.tech/stonework => ./submodule/stonework
   ```

2. Add [CNF Registry][cnf-registry-plugin] and [Punt Manager][punt-manager-plugin] into the list
   of plugins to load by the CNF agent. Some dependencies of these plugins are not injected by default
   and have to be set explicitly as shown below. For example, the CNF Registry requires `CnfIndex` - a CNF
   integer identifier unique among all CNFs (that might be deployed alongside the same StoneWork instance).
   For PANTHEON.tech CNFs, table with CnfIndex assignments can be found [here][cnf-index].
   ```go
   package app

   import (
       cnfreg_plugin "go.pantheon.tech/stonework/plugins/cnfreg"
       puntmgr_plugin "go.pantheon.tech/stonework/plugins/puntmgr"
       "go.pantheon.tech/stonework/proto/cnfreg"
   )
   
   // Index assigned to CNF.
   // Has to be unique among all CNFs that will be deployed with the same StoneWork instance.
   const CnfIndex = 19

   type CnfAgent struct {
      // VPP, Linux and other plugins...
   
      // Plugins for interoperability with StoneWork.
      CnfRegistry *cnfreg_plugin.Plugin
      PuntManager *puntmgr_plugin.Plugin
   }

   func New() *CnfAgent {
       cnfreg_plugin.DefaultPlugin.PuntMgr = &puntmgr_plugin.DefaultPlugin
       cnfreg_plugin.DefaultPlugin.HTTPPlugin = &rest.DefaultPlugin
       cnfreg_plugin.DefaultPlugin.CnfIndex = CnfIndex // Set CNF index (a unique CNF identifier)

       // Init also other plugins (watchers, Orchestrator, etc.)...
   
       switch cnfreg_plugin.DefaultPlugin.GetCnfMode() {
           case cnfreg.CnfMode_STANDALONE:
               // Inject also dependencies of VPP plugins in this case...
           case cnfreg.CnfMode_STONEWORK_MODULE:
               // Disable VPP plugins (VPP is managed by StoneWork)...
               // Plugin is effectively disabled if it is injected as nil.
               puntmgr_plugin.DefaultPlugin.IfPlugin = nil
           case cnfreg.CnfMode_STONEWORK:
               panic("invalid CNF mode")
       }

       return &CnfAgent{
           CnfRegistry:  &cnfreg_plugin.DefaultPlugin,
           PuntManager:  &puntmgr_plugin.DefaultPlugin,
           // Inject also other plugins...
           // VPP plugins should be injected only if CNF runs in the Standalone mode,
           // otherwise leave them disabled (nil references). 
       }
   }  
   ```

3. When CNF runs in the Standalone mode (without StoneWork), it should run its own instance
   of VPP inside the container and manage it by its own CNF agent (i.e. VPP plugins for VPP features
   that are being used have to be initialized by CNF agent). Conversely, CNF running alongside
   StoneWork should not start another instance of VPP (to save resources) and therefore the CNF
   agent should not initialize any VPP plugins (or else `govppmux` plugin fails to connect
   to VPP and the agent will terminate). Using `cnfreg_plugin.DefaultPlugin.GetCnfMode()`
   it is possible to determine the mode at which CNF was started (as shown by the code snippet above).
   This information is determined by the environment variable `CNF_MODE`, which has either
   value `STANDALONE` (default if variable not defined) or `STONEWORK_MODULE`.
   See section `Deployment` from the top-level README.md of StoneWork to learn how to set
   the variable.

4. In order to use the same CNF image regardless of the mode at which it is deployed,
   it is recommended to create two separate config directories - one for the Standalone
   mode and one for the StoneWork-module mode. The most obvious difference in configuration
   is in `supervisor.conf`, which should not include VPP entry unless CNF is running
   in the Standalone mode. Also `initfileregistry.conf` should not enable init-file as
   a configuration source when all input configuration is submitted over StoneWork.
   To select config directory based on the CNF mode, use the following CMD for docker image:
   ```
   CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
       mkdir -p /run/vpp /run/stonework/vpp && \
       if [ "$CNF_MODE" = "STONEWORK_MODULE" ]; then CONFIG_DIR="/etc/cnf-novpp/"; else CONFIG_DIR="/etc/cnf/"; fi && \
       export CONFIG_DIR && \
       exec cnf-init
   ```
   (instead of `/etc/cnf` and `/etc/cnf-novpp` use something more descriptive, e.g. `/etc/dhcp` and `/etc/dhcp-novpp`)

5. CNF Plugins (implementing CRUD operations over CNF config models) will have to depend on the CNF Registry
   plugin and potentially even on the Punt Manager if some packets need to be punted between VPP and CNF/Linux
   (over memif or TAP). Dependency on the CNF Registry is due to a requirement to register all config models
   implemented by the CNF using the method `RegisterCnfModel` as shown below:
   ```go
   type Plugin struct {
       Deps
   }

   type Deps struct {
       infra.PluginDeps
       PuntManager puntmgr_plugin.PuntManagerAPI
       CnfRegistry cnfreg_plugin.CnfAPI
   }
   
   func (p *Plugin) Init() (err error) {
       // init and register descriptors...
   
       // register the model implemented by CNF
       err = p.CnfRegistry.RegisterCnfModel(cnfproto.ModelCnf, cnfDescriptor,
           &cnfreg_plugin.CnfModelCallbacks{
               PuntRequests:     descriptor.CnfPuntReqs,
               ItemDependencies: descriptor.CnfItemDeps,
           })
       if err != nil {
           return err
       }
   }
   ```
   `RegisterCnfModel` takes reference to the model, its descriptor and optionally also callbacks that define
   dependencies and requirements for packet punting. More information can be found in the API interfaces
   `PuntManagerAPI` and `CnfAPI`.

6. Descriptor corresponding to a CNF model will have to behave slightly differently based on the mode
   at which the CNF is deployed (Standalone vs. StoneWork-module).\
   In Standalone mode:
    - If packet punting is required, `Create`/`Delete` methods should call `AddPunt`/`DelPunt` methods of Punt Manager.
      However, interconnect (`memif-memif` or `TAP-TAP`) for punting will not be created just yet, only 
      transaction will be prepared and scheduled for execution. Any operation that depends on the interconnection
      to be prepared should be represented by a derived value that depends on a SB notification published
      by the Punt Manager (key returned by `NotificationKey` from Punt Manager package).
      On the other hand, the metadata for interconnection that will be configured for punting
      are already available immediately after returning from `AddPunt` and can be obtained via `GetPuntMetadata`
      of Punt Manager (e.g. memif/TAP interface names, interface IP/MAC addresses, etc.)
    - If packet punting is required, `Dependencies` should return dependencies determined by `GetPuntDependencies`
      of Punt Manager

   In StoneWork-module mode:
    - If packet punting was required (in `RegisterCnfModel`), at the moment when `Create` is called the punting
      is already established, and `Create` can call `GetPuntMetadata` of Punt Manager to learn metadata about
      the interconnection configured for punting (e.g. memif/TAP interface names, interface IP/MAC addresses, etc.)

7. Additionally to `CNF_MODE` env. variable that was already discussed, one may also need to define variable
   `CNF_MGMT_INTERFACE` or `CNF_MGMT_SUBNET` - these are used by CNF Registry to determine which network interface
   to use to talk to StoneWork (i.e. management interface). `CNF_MGMT_INTERFACE` has higher priority and can be used
   to enter the management interface (host) name directly. `CNF_MGMT_SUBNET` can be used to inform CNF Registry what
   network subnet is used by the management network. Based on that the plugin will be able to determine which interface
   is inside the mgmt network.

8. Another deployment requirement is to mount `/run/stonework/` between StoneWork and every CNF.
   Sub-directory `/run/stonework/discovery` will be created and used by StoneWork and CNFs to discover each other.
   Sub-directory `/run/stonework/memif` will be created and used by Punt Manager for memif sockets.

9. Last requirement targeted for CNF docker image, is that it should have `/api` directory with all the proto files
   that define CNF models (do not include ligato models from upstream, even if used). Directory structure inside `/api`
   should be the same as in the repository (under the `/proto` directory).
   Consider using these commands:
   ```
   RUN mkdir /api
   RUN rsync -v --recursive --chmod=D2775,F444 --exclude '*.go'  proto/ /api/
   ```
   Lastly, api directory should contain file `/api/models.spec.yaml` that contains `ModelDetail` of every
   CNF model (one yaml document per model, separated by "---" with newlines).
   Content of this file can be generated using the `pkg/printspec` package from StoneWork repository.
   Since your CNF already probably has an init process (~supervisor), consider extending it to:
   ```go
   package main

   import (
       "fmt"
       "os"

       "github.com/namsral/flag"

       "go.ligato.io/cn-infra/v2/agent"
       sv "go.ligato.io/cn-infra/v2/exec/supervisor"

       "go.pantheon.tech/stonework/pkg/printspec"

       // include all configuration models exposed by CNF, e.g.:
       "pantheon.tech/my-cnf/proto/cnfproto"
   )

   var printSpec = flag.CommandLine.Bool("print-spec", false,
       "only print spec of CNF models into stdout and exit")

   func main() {
       a := agent.NewAgent(agent.AllPlugins(&sv.DefaultPlugin))
       if *printSpec {
           // List all CNF models in the call to SelectedModels
           err := printspec.SelectedModels(os.Stdout, cnfproto.CnfModel)
           if err != nil {
               _, _ = fmt.Fprint(os.Stderr, err)
               os.Exit(1)
           }
           os.Exit(0)
       }
       if err := a.Run(); err != nil {
           panic(err)
       }
   }
   ```
   Then inside dockerfile call (replace `cnf-init` with the name of your init command):
   ```
   RUN /usr/local/bin/cnf-init --print-spec > /api/models.spec.yaml
   ```
   This api folder can be then used by `submodule/stonework/scripts/gen-docs.sh` to generate
   documentation for the CNF (as markdown and pdf and also with json schema). Script takes two
   arguments: output directory for the documentation (use `docs` in the CNF repo) and the CNF
   name (e.g. "CNF-DHCP"). Script also reads standard input and expect a list of CNF images
   separated by newlines for which a (merged) documentation will be generated.
   The merging of docs can be used to generate documentation for an entire StoneWork
   deployment, that aside from StoneWork itself includes one or more CNFs.
   For your single CNF, generate docs with:
   ```
   $ echo "<your-cnf-image>" | ./submodule/stonework/scripts/gen-docs.sh "./docs/" "<your-cnf-name>"
   ```


*Please note that StoneWork repository contains a reference "mock" CNF implementation
 While it is not a real CNF doing something useful, it certainly acts as one
 and can be deployed either standalone (i.e. running its own copy
 of VPP as data-plane), or it can be discovered and loaded by StoneWork (i.e. share VPP data-plane
 with StoneWork and potentially other CNFs). The mock CNF entry point can be found under
 `cmd/mockcnf`, configuration mock model in `proto/mockcnf/mockcnf.proto` and plugin implementing
 CRUD operations over the model is inside the `plugins/mockcnf` directory.*


[cnf-index]: ../CNF-INDEX.md
[architecture]: ./ARCHITECTURE.md
[punt-manager-plugin]: ../plugins/puntmgr/README.md
[cnf-registry-plugin]: ../plugins/cnfreg/README.md
