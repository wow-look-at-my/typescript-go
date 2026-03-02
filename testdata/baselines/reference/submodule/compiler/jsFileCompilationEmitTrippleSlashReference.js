//// [tests/cases/compiler/jsFileCompilationEmitTrippleSlashReference.ts] ////

//// [a.ts]
class c {
}

//// [b.js]
/// <reference path="c.js"/>
function foo() {
}

//// [c.js]
function bar() {
}


//// [out.js]
"use strict";
class c {
}
function bar() {
}
/// <reference path="c.js"/>
function foo() {
}
