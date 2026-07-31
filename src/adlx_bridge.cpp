#include <iostream>
#include <cstring>

#include "../ADLX/SDK/ADLXHelper/Windows/Cpp/ADLXHelper.h"
#include "../ADLX/SDK/Include/IPerformanceMonitoring.h"

using namespace adlx;

static ADLXHelper helper;
static IADLXSystem* sys;
static IADLXPerformanceMonitoringServicesPtr perf;

extern "C" {
    int init() {
        if (ADLX_FAILED(helper.Initialize())) {
            std::cerr << "ADLX init failed\n";
            return 1;
        }

        if (!(sys = helper.GetSystemServices())) {
            std::cerr << "Failed to get system services\n";
            return 2;
        }

        if (ADLX_FAILED(sys->GetPerformanceMonitoringServices(&perf))) {
            std::cerr << "Failed to get performance monitoring services\n";
            return 3;
        }

        return 0;
    }

    int gpu_info(char* name_buf, int buf_size, int* out_temp, int* out_fan_speed) {
        IADLXGPUListPtr gpu_list;
        if (ADLX_FAILED(sys->GetGPUs(&gpu_list))) {
            std::cerr << "Failed to get GPU list\n";
            return 1;
        }

        IADLXGPUPtr gpu;
        if (ADLX_FAILED(gpu_list->At(0, &gpu))) {
            std::cerr << "Failed to get first GPU\n";
            return 2;
        }

        const char* name = nullptr;
        if (ADLX_SUCCEEDED(gpu->Name(&name))) {
            strncpy(name_buf, name, buf_size - 1);
            name_buf[buf_size - 1] = '\0';
        } else {
            name_buf[0] = '\0';
        }

        IADLXGPUMetricsPtr metrics;
        if (ADLX_FAILED(perf->GetCurrentGPUMetrics(gpu, &metrics))) {
            std::cerr << "Failed to get GPU metrics\n";
            return 3;
        }

        adlx_double temp = 0;
        metrics->GPUTemperature(&temp);
        *out_temp = (int)temp;
        metrics->GPUFanSpeed(out_fan_speed);

        return 0;
    }
}