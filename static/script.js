document.addEventListener("DOMContentLoaded", function(){
    let form = document.getElementById("signup-form");

    form.addEventListener("submit", function(e){
        e.preventDefault();
    
        const payload = new FormData(this);
        console.log(payload);
        fetch("/register", {
            method: 'POST',
            body: payload
        })
        .then(response=>response.json())
        .then(data=>{
            if (data){
                alert("Username already exists.")
            }
        })
    })
})