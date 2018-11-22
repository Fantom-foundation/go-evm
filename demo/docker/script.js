var Eth = require('web3-eth');

//var eth = new Eth(Eth.givenProvider || 'ws://evm:8546');
var eth = new Eth(Eth.givenProvider || 'http://evm:8545');


console.log("01 getProtocolVersion");
eth.getProtocolVersion().then(console.log, console.log);

//console.log("02 getAccounts");
//eth.getAccounts(console.log);


//console.log("03 sendTransaction");
//eth.sendTransaction({
//    from: '0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe',
//    to: '0x11f4d0A3c12e86B4b5F39B213F7E19D048276DAe',
//    value: '1000000000000000'
//})
//.then(console.log, console.log);


console.log("Finish");
