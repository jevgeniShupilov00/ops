package main

import (
    "fmt"
    "os"
    "os/exec"
)

type drive struct {
    path   string
    format string
}

type netdev struct {
    nettype string
    id      string
    ifname  string
    mac     string
}

type qemu struct {
    drives []drive
    ifaces []netdev
}

func (q *qemu) addDrive(path string, format string) {

}

func (q *qemu) addNetworkDevice(n netdev) {

}

func (q *qemu) start(image string , port int) error {
    // TODO: https://github.com/deferpanic/uniboot/issues/31	
    err := copy(finalImg,"image2");	
    panicOnError(err)
    args := q.Args(port)
    fmt.Println(args)
    cmd := exec.Command("qemu-system-x86_64", args...)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Println(err)
        return err
    }
    return nil
}

func (q *qemu) Args(port int) []string {
    // TODO : this should come from q.drives and ifaces
    args := []string{}

    boot := []string{"-drive", "file=image,format=raw,index=0"}
    storage := []string{"-drive", "file=image2,format=raw,if=virtio"}
    var net []string
    if port > 0 {
        // hostfwd=tcp::8080-:8080
        portfw := fmt.Sprintf("hostfwd=tcp::%v-:%v", port, port)
        net = []string{"-device", "virtio-net,netdev=n0","-netdev", "user,id=n0," + portfw}
      
    } else {
        net = []string{"-device", "virtio-net,mac=7e:b8:7e:87:4a:ea,netdev=n0",
        "-netdev", "tap,id=n0,ifname=tap0,script=no,downscript=no"}
    }
    display := []string{"-display", "none", "-serial", "stdio"}
    args = append(args, boot...)
    args = append(args, display...)
    args = append(args, []string{"-nodefaults", "-no-reboot", "-m", "2G", "-device", "isa-debug-exit"}...)
    args = append(args, storage...)
    args = append(args, net...)
    args = append(args, "-enable-kvm")
    return args
}

func newQemu() Hypervisor {
    return &qemu{}
}