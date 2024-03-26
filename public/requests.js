async function postData(url = "", data = "") {
    const response = await fetch(url, {
        method: "POST",
        credentials: "same-origin",
        headers: {
            "Content-Type": "text/plain",
        },
        redirect: "follow",
        referrerPolicy: "no-referrer",
        body: data,
    });
    if (response.ok) {
        return response.text();
    } else {
        return response.status
    }
}

async function getData(url = "", data = "") {
    const response = await fetch(url, {
        method: "GET",
        credentials: "same-origin",
        redirect: "follow",
        referrerPolicy: "no-referrer"
    });
    if (response.ok) {
        return response.text();
    } else {
        console.log(response.text())
        return response.status
        // return "Error"
    }
}

function saveLog() {
    let logField = document.getElementById("log")
    postData("/api/log", logField.value).then((data) => {
        if (data !== 500) {
            logField.value = ""
        }
    });
}

function saveToday() {
    let logField = document.getElementById("day")
    postData("/api/today", logField.value).then((data) => {
        console.log(data);
    });
}
function loadToday() {
    let dayField = document.getElementById("day")
    getData("/api/today", dayField.value).then((data) => {
        if (data === 404) {
            dayField.value = ""
        } else if (data === 401) {
            dayField.enabled = false
            dayField.value = "Unauthorized"
        } else if (data === 500) {
            dayField.enabled = false
            dayField.value = "Internal server error"
        } else {
            dayField.value = data
        }
    });
}

function loadPrevious() {
    let daysField = document.getElementById("days")
    daysField.innerHTML = ""
    getData("/api/day", "").then((data) => {
        if (data === 401) {
            alert("Unauthorized")
        } else if (data === 500) {
            alert("Internal server error")
        } else {
            data = JSON.parse(data).reverse()  // Reverse: latest days first
            for (let i in data) {
                let li = document.createElement("li");
                li.innerHTML = `<a href="/api/day/${data[i]}">${data[i]}</a>`
                daysField.appendChild(li);
            }
        }
    });
}