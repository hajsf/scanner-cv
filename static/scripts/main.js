function exportTableToExcel(tableID, filename = ''){
    var downloadLink;
    var dataType = 'application/vnd.ms-excel';
    var tableSelect = document.getElementById(tableID);
    var tableHTML = tableSelect.outerHTML.replace(/ /g, '%20');
    
    // Specify file name
    filename = filename?filename+'.xls':'excel_data.xls';
    
    // Create download link element
    downloadLink = document.createElement("a");
    
    document.body.appendChild(downloadLink);
    
    if(navigator.msSaveOrOpenBlob){
        var blob = new Blob(['\ufeff', tableHTML], {
            type: dataType
        });
        navigator.msSaveOrOpenBlob( blob, filename);
    }else{
        // Create a link to the file
        downloadLink.href = 'data:' + dataType + ', ' + tableHTML;
    
        // Setting the file name
        downloadLink.download = filename;
        
        //triggering the function
        downloadLink.click();
    }
}

async function send(params) {
    var path = document.querySelector('#path');
    var skills = document.querySelector('#skills');
    var skillsList = skills.value
    skillsList = skillsList.replace(/\n\r?/g, '|');
    const dataToSend = JSON.stringify({"path": path.value, "skills": skillsList});
        fetch("http://localhost:3000/getSkills", {
            credentials: "same-origin",
            mode: "cors",
            method: "post",
            headers: { "Content-Type": "text/plain; charset=utf-8"}, //"application/json" },
            body: dataToSend
        }).then(response => {
            if (response.status === 200) {
                return getTextFromStream(response.body)
            } else {
                console.log("Status: " + response.status)
                return Promise.reject("server")
            }
        }).then(dataText => {
               console.log(`dataJson: ${dataText}`)
           })
           .catch(err => {
               if (err === "server") return
               console.log(err)
           })    
}
/*
function send(params) {
    var path = document.querySelector('#path');
    var skills = document.querySelector('#skills');
    var skillsList = skills.value
    skillsList = skillsList.replace(/\n\r?/g, '|');
    const dataToSend = JSON.stringify({"path": path.value, "skills": skillsList});
    let dataReceived = ""; 
    fetch("http://localhost:3000/getSkills", {
        credentials: "same-origin",
        mode: "cors",
        method: "post",
        headers: { "Content-Type": "text/plain; charset=utf-8"}, //"application/json" },
        body: dataToSend
    }).then(resp => {
        if (resp.status === 200) {
            return resp.json();
        } else {
            console.log("Status: " + resp.status)
            return Promise.reject("server")
        }
    }).then(dataJson => {
        dataReceived = JSON.parse(dataJson)
        console.log(`Received: ${dataReceived}`)
    }).catch(err => {
        if (err === "server") return
        console.log(err)
    })
} */