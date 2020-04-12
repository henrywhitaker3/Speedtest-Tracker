dojoConfig = {
  async: true,
  baseUrl: "lib_dojo/",
  packages: [{
    name: "dojo",
    location: "dojo"
  },{
    name: "dojox",
    location: "dojox"
  }],
  deps: ["export"],
  highlightJsDir: __dirname + "/lib_highlight/build/lib",
  cwd: __dirname,
};

require("./lib_dojo/dojo/dojo.js");
