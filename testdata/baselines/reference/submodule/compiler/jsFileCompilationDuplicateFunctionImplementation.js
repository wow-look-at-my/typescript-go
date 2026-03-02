//// [tests/cases/compiler/jsFileCompilationDuplicateFunctionImplementation.ts] ////

//// [b.js]
function foo() {
    return 10;
}
//// [a.ts]
function foo() {
    return 30;
}



//// [out.js]
"use strict";
function foo() {
    return 10;
}
function foo() {
    return 30;
}
