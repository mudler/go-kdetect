# go-kdetect

A small and simple library to report kernel modules that are detected in a running system.

It uses [go-udev](github.com/pilebones/go-udev) to detect the system devices and parses the kernel alias file to derive the drivers to load.

## Usage


```golang
package main


import (
    "github.com/mudler/go-kdetect"
    "fmt"
)

func main() {
    drivers, err := kdetect.ProbeKernelModules("/lib/modules/5.11.15-1-default/modules.alias")
    if err != nil {
        // Handle err
    }

    // drivers is a []string
    fmt.Println(drivers)
    // [acpi_cpufreq pcc_cpufreq tiny_power_button button ac processor_aggregator int3402_thermal video fjes dell_smo8800 smo8800 int3403_thermal ec battery thermal skl_uncore nouveau pcieport i915 proc_thermal xhci_hcd pcmcia hub hid_generic logitech-djreceiver usbhid hid_logitech_hidpp logitech-hidpp-device usb uvcvideo btusb intel_pch_thermal des_generic i2c_designware idma64 intel-lpss hid-multitouch i2c_hid mei_wdt mei_hdcp mei_me serial sr_mod sd ahci iwlwifi rtsx_pci_sdmmc rtsx_pci int3403 thermal snd_hda_codec_realtek snd_soc_hdac_hdmi snd_hda_codec_hdmi snd_hda_intel dummy ee1004 iTCO_wdt i801_smbus int3402 thermal acpi-wmi mxm_wmi intel_pmc_core intel-hid int3400 thermal tpm_tis acpi-fan intel_wmi_thunderbolt intel-wmi-thunderbolt wmi_bmof wmi-bmof dell_wmi_descriptor dell-wmi-descriptor dell_wmi dell-wmi dell_smbios dell-smbios coretemp dcdbas dell-laptop efi-framebuffer efi_pstore atkbd psmouse i8042 intel_rapl_msr kgdboc pcspkr reg-dummy serial8250 snd-soc-dummy vboxdrv system alarmtimer rtc_cmos i8042 kbd i8042 aux intel_rapl_common processor]

}
```
