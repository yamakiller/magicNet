#include "_cgo_export.h"
#include <stdint.h>
#include <stdio.h>
#include <string.h>

int go_output_wrapper(const char *buf, int len,  ikcpcb *kcp, void *user) {   
    GoSlice slice;
    slice.data = (void*)buf;
    slice.len  = len;
    slice.cap  = len;
    return go_output(slice, (GoUintptr)user);
}


ikcpcb* mkcp_create(IUINT32 conv, uintptr_t user) {
    ikcpcb *kcp = ikcp_create(conv, (void*)user);
    kcp->output = go_output_wrapper;
    return kcp;
}

void mkcp_release(ikcpcb *kcp) {
    return ikcp_release(kcp);
}

IUINT32 mkcp_getconv(char* ptr) {
    return ikcp_getconv((const void*)ptr);
}

int mkcp_send(ikcpcb *kcp, char *buffer, int len) {
    return ikcp_send(kcp,(const char*)buffer, len);
}


int mkcp_input(ikcpcb *kcp, char *buffer, int len) {
    return ikcp_input(kcp, (const char*)buffer, len);
}