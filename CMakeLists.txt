cmake_minimum_required(VERSION 3.30)

# CXX aka CPP
project("live-subtitles" CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED True)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON) # Set for better clangd integrations

include(FindPkgConfig)

file(GLOB_RECURSE MY_SOURCES CONFIGURE_DEPENDS src/*.cpp)
add_executable(${PROJECT_NAME} ${MY_SOURCES})

# Dependencies
pkg_check_modules(GST REQUIRED gstreamer-1.0)
target_include_directories(${PROJECT_NAME} PUBLIC ${GST_INCLUDE_DIRS})
target_compile_options(${PROJECT_NAME} PUBLIC ${GST_CFLAGS_OTHER})
target_link_libraries(${PROJECT_NAME} ${GST_LIBRARIES})

pkg_check_modules(CURL REQUIRED curlpp)
target_include_directories(${PROJECT_NAME} PUBLIC ${CURL_INCLUDE_DIRS})
target_compile_options(${PROJECT_NAME} PUBLIC ${CURL_CFLAGS_OTHER})
target_link_libraries(${PROJECT_NAME} ${CURL_LIBRARIES})

pkg_check_modules(AWS REQUIRED aws-cpp-sdk-core)
target_include_directories(${PROJECT_NAME} PUBLIC ${AWS_INCLUDE_DIRS})
target_compile_options(${PROJECT_NAME} PUBLIC ${AWS_CFLAGS_OTHER})
target_link_libraries(${PROJECT_NAME} ${AWS_LIBRARIES})

pkg_check_modules(PORTAUDIO REQUIRED portaudio-2.0)
target_include_directories(${PROJECT_NAME} PUBLIC ${PORTAUDIO_INCLUDE_DIRS})
target_compile_options(${PROJECT_NAME} PUBLIC ${PORTAUDIO_CFLAGS_OTHER})
target_link_libraries(${PROJECT_NAME} ${PORTAUDIO_LIBRARIES})

target_compile_options(${PROJECT_NAME} PRIVATE -fexceptions)
