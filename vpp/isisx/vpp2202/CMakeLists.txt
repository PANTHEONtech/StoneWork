cmake_minimum_required(VERSION 3.5)

project(isisx-plugin)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

find_package(VPP)

add_subdirectory(isisx)
