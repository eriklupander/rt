//+build !noasm !appengine
// AUTO-GENERATED BY C2GOASM -- DO NOT EDIT

TEXT ·__MultiplyMatrixByVec64(SB), $0-24

	MOVQ m+0(FP), DI
	MOVQ vec4+8(FP), SI
	MOVQ result+16(FP), DX

	LONG $0x0610fdc5               // vmovupd    ymm0, yword [rsi]
	LONG $0x0f59fdc5               // vmulpd    ymm1, ymm0, yword [rdi]
	LONG $0x5759fdc5; BYTE $0x20   // vmulpd    ymm2, ymm0, yword [rdi + 32]
	LONG $0x5f59fdc5; BYTE $0x40   // vmulpd    ymm3, ymm0, yword [rdi + 64]
	LONG $0x4759fdc5; BYTE $0x60   // vmulpd    ymm0, ymm0, yword [rdi + 96]
	LONG $0x197de3c4; WORD $0x01cc // vextractf128    xmm4, ymm1, 1
	LONG $0xcc58f1c5               // vaddpd    xmm1, xmm1, xmm4
	LONG $0x197de3c4; WORD $0x01d4 // vextractf128    xmm4, ymm2, 1
	LONG $0xd458e9c5               // vaddpd    xmm2, xmm2, xmm4
	LONG $0x197de3c4; WORD $0x01dc // vextractf128    xmm4, ymm3, 1
	LONG $0xdc58e1c5               // vaddpd    xmm3, xmm3, xmm4
	LONG $0x197de3c4; WORD $0x01c4 // vextractf128    xmm4, ymm0, 1
	LONG $0xc458f9c5               // vaddpd    xmm0, xmm0, xmm4
	LONG $0x1875e3c4; WORD $0x01ca // vinsertf128    ymm1, ymm1, xmm2, 1
	LONG $0xd314f5c5               // vunpcklpd    ymm2, ymm1, ymm3
	LONG $0x01fde3c4; WORD $0xd8d2 // vpermpd    ymm2, ymm2, 216
	LONG $0x197de2c4; BYTE $0xe0   // vbroadcastsd    ymm4, xmm0
	LONG $0x0d6de3c4; WORD $0x08d4 // vblendpd    ymm2, ymm2, ymm4, 8
	LONG $0xcb15f5c5               // vunpckhpd    ymm1, ymm1, ymm3
	LONG $0x01fde3c4; WORD $0xd8c9 // vpermpd    ymm1, ymm1, 216
	LONG $0x187de3c4; WORD $0x01c0 // vinsertf128    ymm0, ymm0, xmm0, 1
	LONG $0x0d75e3c4; WORD $0x08c0 // vblendpd    ymm0, ymm1, ymm0, 8
	LONG $0xc058edc5               // vaddpd    ymm0, ymm2, ymm0
	LONG $0x0211fdc5               // vmovupd    yword [rdx], ymm0
	VZEROUPPER
	RET