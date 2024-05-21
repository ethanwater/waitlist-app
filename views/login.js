function strip(s) {
    return s.replace(/\s+/g, '');
}

async function enrollEmailResponse(endpoint, email, password, aborter) {
    const response = await fetch(`/${endpoint}?q=${encodeURIComponent(email)}`, {
        signal: aborter.signal,
    });
    const text = await response.text();

    if (response.ok) {
        return text;
    } else {
        throw new Error(text);
    }
}

async function fetchKey(endpoint, aborter) {
    const response = await fetch(`/login/${endpoint}`, {
        signal: aborter
    });
    const text = await response.text();
    if (response.ok) {
        return text;
    } else {
        throw new Error(text);
    }
}

async function verifyCode(endpoint, hash, input, aborter) {
    const response = await fetch(`/login/${endpoint}?hash=${encodeURIComponent(hash)}&input=${encodeURIComponent(input)}`, {
        signal: aborter.signal,
    });
    const text = await response.text();

    if (response.ok) {
        console.log("authentication code success")
        return text;
    } else {
        console.log("authentication code failure")
        throw new Error(text);
    }
}

async function retrieveCode() {
    var inputValues = Array.from(document.querySelectorAll('.symbol')).map(input => input.value);
    var input = inputValues.toString().replaceAll(",", "");
    const handleIncorrectAnimation = () => {
        document.getElementById("verify").classList.remove("incorrect");
    };

    if (strip(input).length < 5) {
        return none;
    }

    var controller = new AbortController();
    for (const endpoint of['verifykey']) {
        const responseText = await verifyCode(endpoint, getCookie2FA("vivian2FA"), input.toUpperCase(), controller);
        const results = JSON.parse(responseText);

        if (!results) {
            document.getElementById("verify").classList.add("incorrect");
            document.getElementById("verify").addEventListener("animationend", handleIncorrectAnimation, {
                once: true
            });
            errorMessage("incorrect code");
            console.log(results);
        } else {
            if (document.getElementById('error')) {
                document.getElementById('error').style.visibility = 'false';
                document.getElementById('error').style.display = 'none';
            }
            console.log(results);
            window.location.assign("../apps-chart/index.html")
        }
    }
}

function createVerificationElement() {
    if (!document.getElementById("verifyid")) {
        var container = document.getElementById("verify");
        var verificationDiv = document.createElement("div");
        verificationDiv.className = "verification";
        verificationDiv.id = "verifyid";

        for (var i = 1; i <= 5; i++) {
            var input = document.createElement("input");
            input.id = "code" + i;
            input.type = "code";
            input.className = "symbol";
            input.setAttribute("oninput", "focusNextInput(this, event)");
            input.required = true;

            verificationDiv.appendChild(input);
        }

        container.appendChild(verificationDiv);
        document.getElementById("code1").focus();
    }
}

function errorMessage(msg) {
    if (!document.getElementById("error")) {
        var div = document.createElement("div");
        div.id = "error"
        div.className = "errormsg";
        div.innerText = msg;

        document.getElementById("main").appendChild(div);
    } else {
        var div = document.getElementById("error");
        div.innerText = msg;
    }
}

function main() {
    const email = document.getElementById('email');
    const enterButton = document.getElementById('enter');
    const inputs = document.querySelectorAll('input');

    inputs.forEach((input) => {
        input.setAttribute('autocomplete', 'off');
        input.setAttribute('autocorrect', 'off');
        input.setAttribute('autocapitalize', 'off');
        input.setAttribute('spellcheck', false);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    main(); 
});
