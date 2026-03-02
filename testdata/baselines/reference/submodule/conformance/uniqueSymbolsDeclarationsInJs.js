//// [tests/cases/conformance/types/uniqueSymbol/uniqueSymbolsDeclarationsInJs.ts] ////

//// [uniqueSymbolsDeclarationsInJs.js]
// classes
class C {
    /**
     * @readonly
     */
    static readonlyStaticCall = Symbol();
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyStaticType;
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyStaticTypeAndCall = Symbol();
    static readwriteStaticCall = Symbol();

    /**
     * @readonly
     */
    readonlyCall = Symbol();
    readwriteCall = Symbol();
}

/** @type {unique symbol} */
const a = Symbol();


//// [uniqueSymbolsDeclarationsInJs-out.js]
"use strict";
// classes
class C {
    /**
     * @readonly
     */
    static readonlyStaticCall = Symbol();
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyStaticType;
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyStaticTypeAndCall = Symbol();
    static readwriteStaticCall = Symbol();
    /**
     * @readonly
     */
    readonlyCall = Symbol();
    readwriteCall = Symbol();
}
/** @type {unique symbol} */
const a = Symbol();
