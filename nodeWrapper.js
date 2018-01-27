const childProcess = require("child_process").execFile;
const execPath = "./goBundler.exe";

const configFileName = process.argv[2] || "";
const params = ["config.json", configFileName];

const child = childProcess(execPath, params);
child.stdout.on("data", data => {
  console.log(data.slice(0, -1));
});
