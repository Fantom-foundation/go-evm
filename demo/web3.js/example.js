// https://web3js.readthedocs.io/en/1.0/web3-eth.html

var Eth = require('web3-eth');

var eth = new Eth(Eth.givenProvider || 'http://localhost:8545');

console.log("01 Create wallet");
//var wallet = eth.accounts.wallet.create();
console.log(eth.accounts.wallet);

console.log("02 Create account 1");
var acc1 = eth.accounts.create();
console.log(acc1);

console.log("03 Create account 2");
var acc2 = eth.accounts.create(0);
console.log(acc2);

console.log("04 Save wallet");
eth.accounts.wallet.add(acc1);
eth.accounts.wallet.add(acc2);
//eth.accounts.wallet.save("myPassword99");

//console.log(wallet);
//console.log(eth.accounts.wallet);

console.log("04 getAccounts");
eth.getAccounts(console.log);

console.log("05 sendTransaction");
eth.sendTransaction({
    from: acc1.address,
    to: acc2.address,
    value: '1000000000000000'
})
.then(console.log, console.log)
