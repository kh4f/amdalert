#![allow(non_camel_case_types)]
use libloading::{Library, Symbol};
use std::ffi::c_void;

#[repr(C)]
struct ADLTemperature { size: i32, temperature: i32 }
unsafe extern "C" fn adl_malloc(size: i32) -> *mut c_void { unsafe { libc::malloc(size as usize) } }
type ADL_MAIN_CONTROL_CREATE = unsafe extern "C" fn(unsafe extern "C" fn(i32) -> *mut c_void, i32) -> i32;
type ADL_MAIN_CONTROL_DESTROY = unsafe extern "C" fn() -> i32;
type ADL_ADAPTER_NUMBER_OF_ADAPTERS_GET = unsafe extern "C" fn(*mut i32) -> i32;
type ADL_OVERDRIVE5_TEMPERATURE_GET = unsafe extern "C" fn(i32, i32, *mut ADLTemperature) -> i32;

struct ADLManager<'a> { adl_destroy: Symbol<'a, ADL_MAIN_CONTROL_DESTROY> }
impl Drop for ADLManager<'_> {
    fn drop(&mut self) {
        unsafe { (*self.adl_destroy)(); }
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    unsafe {
        let lib = Library::new("atiadlxx.dll").expect("Failed to load atiadlxx.dll");

        let adl_create: Symbol<ADL_MAIN_CONTROL_CREATE> = lib.get(b"ADL_Main_Control_Create")?;
        let adl_destroy: Symbol<ADL_MAIN_CONTROL_DESTROY> = lib.get(b"ADL_Main_Control_Destroy")?;
        let get_adapters: Symbol<ADL_ADAPTER_NUMBER_OF_ADAPTERS_GET> = lib.get(b"ADL_Adapter_NumberOfAdapters_Get")?;
        let get_temp: Symbol<ADL_OVERDRIVE5_TEMPERATURE_GET> = lib.get(b"ADL_Overdrive5_Temperature_Get")?;

        adl_create(adl_malloc, 1);
		let _manager = ADLManager { adl_destroy };

        let mut adapters = 0;
        get_adapters(&mut adapters);
        if adapters == 0 { return Err("No adapters found".into()) }
        println!("Adapters found: {adapters}");

        let mut temp = ADLTemperature {
    		size: std::mem::size_of::<ADLTemperature>() as i32,
			temperature: 0,
		};
		match get_temp(0, 0, &mut temp) {
			0 => println!("GPU Temperature: {}°C", temp.temperature / 1000),
			_ => println!("Temperature read failed"),
		}

        Ok(())
    }
}