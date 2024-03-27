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
    return response
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
        if (data.ok) {
            logField.value = ""
        }
    });
}

function saveToday() {
    let logField = document.getElementById("day")
    postData("/api/today", logField.value).then((data) => {
        if (!data.ok) {
            alert(`Error saving: ${data.text()}`)
        }
    });
}
function loadToday() {
    let dayField = document.getElementById("day")
    getData("/api/today", dayField.value).then((data) => {
        if (data === 404) {
            dayField.value = ""
        } else if (data === 401) {
            dayField.readOnly = true
            dayField.value = "Unauthorized"
        } else if (data === 500) {
            dayField.readOnly = true
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
                li.innerHTML = `<a href="/day?d=${data[i]}">${formatDate(data[i])}</a>`  // Parse to human-readable
                daysField.appendChild(li);
            }
        }
    });
}

function loadDay() {
    const urlParams = new URLSearchParams(window.location.search);
    const day = urlParams.get('d');

    let dayTag = document.getElementById("daytag")
    dayTag.innerText = formatDate(day)

    let dayField = document.getElementById("day")
    getData("/api/day/" + day, "").then((data) => {
        if (data === 404) {
            dayField.value = ""
        } else if (data === 401) {
            dayField.value = "Unauthorized"
        } else if (data === 500) {
            dayField.value = "Internal server error"
        } else {
            dayField.value = data
        }
    });
    dayField.readOnly = true
}