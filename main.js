let os = require('os');
let child;
let fails = 0;
let goBinary = "./mainnet-commitments-ui";
let args = []

function setPage(html) {
    const container = document.getElementById("app");
    container.innerHTML = html;

    // Set focus for autofocus element
    let elem = document.querySelector("input[autofocus]");
    if (elem != null) {
        elem.focus();
    }
}

function body_message(msg) {
    setPage('<h1>' + msg + '</h1>');
}

function start_process() {
    body_message("Loading...");

    const spawn = require('child_process').spawn;
    child = spawn(goBinary, args);

    const readline = require('readline');
    const rl = readline.createInterface({
        input: child.stdout
    })

    rl.on('line', (data) => {
        console.log(`Received: ${data}`);

        if (data.charAt(0) === "$") {
            data = data.substring(1);
            eval(data);
        } else {
            setPage(data);
        }
    });

    child.stderr.on('data', (data) => {
        console.log(`stderr: ${data}`);
    });

    child.on('close', (code) => {
        body_message(`process exited with code ${code}`);
        restart_process();
    });

    child.on('error', (err) => {
        body_message('Failed to start child process: ' + err);
        restart_process();
    });
}

function restart_process() {
    setTimeout(function () {
        fails++;
        if (fails > 5) {
            close();
        } else {
            start_process();
        }
    }, 5000);
}

function element_as_object(elem) {
    let obj = {
        properties: {}
    }
    for (let j = 0; j < elem.attributes.length; j++) {
        obj.properties[elem.attributes[j].name] = elem.attributes[j].value;
    }
    //overwrite attributes with properties
    if (elem.value != null) {
        obj.properties["value"] = elem.value.toString();
    }
    if (elem.checked != null && elem.checked) {
        obj.properties["checked"] = "true";
    } else {
        delete (obj.properties["checked"]);
    }
    return obj;
}

function element_by_tag_as_array(tag) {
    let items = [];
    let elems = document.getElementsByTagName(tag);
    for (let i = 0; i < elems.length; i++) {
        items.push(element_as_object(elems[i]));
    }
    return items;
}

function fire_event(name, sender) {
    let msg = {
        name: name,
        sender: element_as_object(sender),
        inputs: element_by_tag_as_array("input").concat(element_by_tag_as_array("select"))
    }
    child.stdin.write(JSON.stringify(msg));
    console.log(JSON.stringify(msg));
}

function fire_keyPressed_event(e, keycode, name, sender) {
    if (e.keyCode === keycode) {
        e.preventDefault();
        fire_event(name, sender);
    }
}

function avoid_reload() {
    if (sessionStorage.getItem("loaded") === "true") {
        alert("go-webkit will fail when page reload. avoid using <form> or submit.");
        close();
    }
    sessionStorage.setItem("loaded", "true");
}

if (os.platform().isWindows) {
    goBinary += ".exe";
}

avoid_reload();
start_process();


function changeRangeValue(val, max) {
    val = val.replace(/\D/g, '');
    if (val > max) {
        val = max;
    } else if (val < 0) {
        val = 0;
    } else if (val == "-0") {
        val = 0;
    } else if (val % 1 != 0) {
        val = Math.floor(val);
    } else if (val.length > 7) {
        val = val.slice(0, 7);
    }
    if (val.length > 0) {
        document.getElementById("number").value = val;
    }
    document.getElementById("range").value = isNaN(parseInt(val, 10)) ? 0 : parseInt(val, 10);
}

function changeInputValue(val) {
    document.getElementById("number").value = isNaN(parseInt(val, 10)) ? 0 : parseInt(val, 10);
}