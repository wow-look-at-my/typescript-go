//// [tests/cases/compiler/jsFileCompilationEmitDeclarations.ts] ////

//// [a.ts]
class c {
}

//// [b.js]
function foo() {
}


//// [out.js]
"use strict";
class c {
}
function foo() {
}
