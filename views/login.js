function strip(s) {
    return s.replace(/\s+/g, '');
}

async function sendVerificationCode(email, aborter) {
    try {
        const response = await fetch(`/api/sendverificationcode?email=${email}`, {
            signal: aborter.signal,
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.statusText}`);
        }

        const data = await response.json(); // Parse the JSON response
        return data;
    } catch (error) {
        console.error('Error sending verification code:', error);
        throw error;
    }
}

async function verifyVerificationCode(input, aborter) {
    try {
        var inputCapitalized = input.toUpperCase();
        const response = await fetch(`/api/verifyverificationcode?code=${inputCapitalized}`, {
            signal: aborter.signal,
        });

        const data = await response.json(); 
        return data;
    } catch (error) {
        console.error('Error validating code:', error);
        throw error;
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

function getVerificationValues() {
    var values = [];
    for (var i = 1; i <= 5; i++) {
        var input = document.getElementById("code" + i);
        if (input) {
            values.push(input.value);
        }
    }
    console.log("Verification code values:", values.join(''));
    return values.join('');
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
    const verifyDiv = document.getElementById('verify');
    const result = document.getElementById('result');
    var hashedCode;

    const register = async () => {
        var controller = new AbortController();
        try {
            const response = await sendVerificationCode(email.value, controller);
            hashedCode = response.code;
            createVerificationElement();                
        } catch (error) {
            errorMessage('Failed to send verification email');
        }
    } 

    const verify = async () => {
        var controller = new AbortController();
        var values = getVerificationValues();
        try {
            const response = await verifyVerificationCode(values, controller)
            console.log('Received response:', response); 
            if (response.code == 'true') {
                console.log('verified');
                email.style.display='none';
                verifyDiv.style.display='none';
                enterButton.style.display='none';

                result.innerText="You're in! 🚀";
            } else {
                errorMessage('invalid code')
            }
        } catch(error) {
            errorMessage('invalid code')
        }
    }

    inputs.forEach((input) => {
        input.setAttribute('autocomplete', 'off');
        input.setAttribute('autocorrect', 'off');
        input.setAttribute('autocapitalize', 'off');
        input.setAttribute('spellcheck', false);
    });

    enterButton.addEventListener('click', () => {
        if (!document.getElementById("verifyid")) {
            register();
        } else {
            verify();
        }
    })
}

document.addEventListener('DOMContentLoaded', () => {
    main(); 
});
