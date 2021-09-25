#include "kcp/ikcp.h"
#include <stdint.h>

ikcpcb* mkcp_create(IUINT32 conv, uintptr_t user);

void mkcp_release(ikcpcb *kcp);

IUINT32 mkcp_getconv(char *ptr);

int mkcp_send(ikcpcb *kcp, char *buffer, int len);

int mkcp_input(ikcpcb *kcp, char *buffer, int len);
