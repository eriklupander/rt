//+build !noasm !appengine
// AUTO-GENERATED BY C2GOASM -- DO NOT EDIT

TEXT ·__CrossProduct(SB), $0-24

	MOVQ vec1+0(FP), DI
	MOVQ vec2+8(FP), SI
	MOVQ result+16(FP), DX

	LONG $0x0710fdc5               // vmovupd    ymm0, yword [rdi]
	LONG $0x0e10fdc5               // vmovupd    ymm1, yword [rsi]
	LONG $0x01fde3c4; WORD $0xc9d1 // vpermpd    ymm2, ymm1, 201
	LONG $0xd259fdc5               // vmulpd    ymm2, ymm0, ymm2
	LONG $0x01fde3c4; WORD $0xc9c0 // vpermpd    ymm0, ymm0, 201
	LONG $0xc059f5c5               // vmulpd    ymm0, ymm1, ymm0
	LONG $0xc05cedc5               // vsubpd    ymm0, ymm2, ymm0
	LONG $0x01fde3c4; WORD $0xc9c0 // vpermpd    ymm0, ymm0, 201
	LONG $0x0211fdc5               // vmovupd    yword [rdx], ymm0
	VZEROUPPER
	RET
