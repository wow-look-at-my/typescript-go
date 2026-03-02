//// [tests/cases/compiler/outModuleConcatUnspecifiedModuleKind.ts] ////

//// [a.ts]
export class A { } // module

//// [b.ts]
var x = 0; // global

//// [out.js]
export class A {
} // module
var x = 0; // global
