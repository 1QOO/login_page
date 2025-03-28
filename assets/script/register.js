
document.addEventListener("DOMContentLoaded", function(){

    let form = document.getElementById("signup-form");
    let usernameMessage = document.getElementById("username-message");
    let passwordMessage = document.getElementById("password-message");
    let username = document.getElementById("username");
    let sessionToken = document.getElementById("session-token");
    let password = document.getElementById("password");
    let confirmedPassword = document.getElementById("confirmed-password");

    function ShowPasswordMessage(message){
        password.focus();
        password.value = "";
        confirmedPassword.value = "";
        passwordMessage.innerText = message;
        password.addEventListener("input", 
            function hideMessage(){
                passwordMessage.innerText = "";
                username.removeEventListener("input", hideMessage);
                confirmedPassword.removeEventListener("input", hideMessage);
            }
        )
        confirmedPassword.addEventListener("input", 
            function hideMessage(){
                passwordMessage.innerText = "";
                password.removeEventListener("input", hideMessage);
                confirmedPassword.removeEventListener("input", hideMessage);
            }
        )
    }

    form.addEventListener("submit", function(e){
        console.log("Form default prevented");

        if (password.value.length >= 8){

            if (password.value === confirmedPassword.value){

                const payload = new FormData(this);
                
                fetch("/signup", {
                    method: 'POST',
                    body: payload
                })
                .then(response=>response.json())
                .then(data=>{
                    switch (data.Status){
                        case "UsernameIsTaken" : ShowUsernameMessage(data.Message);
                        case"PasswordDidNotMatch" : ShowPasswordMessage(data.Message);
                        case "PasswordIsTooShort" : ShowPasswordMessage(data.Message);
                        case "Registered" : {
                            window.location.href = data.Message;
                        }
                        default : return
                    
                    }
                })
                .catch(error=>{
                    console.error("Error: ", error);
                })
            }
            else {
                ShowPasswordMessage("Password doesn't match")
            }
        }
        else {
            ShowPasswordMessage("Password must contain 8 or more character");
        }
    })
//END                
})

function ShowUsernameMessage(message){
    usernameMessage.innerText = message;
    username.focus();
    username.value = "";
    username.addEventListener("input", 
        function hideMessage(){
            usernameMessage.innerText = "";
            username.removeEventListener("input", hideMessage);
        }
    )
}

function ShowPasswordMessage(message){
    password.focus();
    password.value = "";
    confirmedPassword.value = "";
    passwordMessage.innerText = message;
    password.addEventListener("input", 
        function hideMessage(){
            passwordMessage.innerText = "";
            username.removeEventListener("input", hideMessage);
            confirmedPassword.removeEventListener("input", hideMessage);
        }
    )
    confirmedPassword.addEventListener("input", 
        function hideMessage(){
            passwordMessage.innerText = "";
            password.removeEventListener("input", hideMessage);
            confirmedPassword.removeEventListener("input", hideMessage);
        }
    )
}