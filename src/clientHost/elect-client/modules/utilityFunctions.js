const fs = require('fs')

// Function that returns text after the html escape characrers have been replaced with html codes
//
// List of all escape characters (in case we want to add more): https://www.werockyourweb.com/html-escape-characters/
//
exports.sanatizeText = (input)=>{

    let output = input;

    if(!output){
        return ""
    }

    output = output.replaceAll('<','&lt;');
    output = output.replaceAll('>','&gt;');
    output = output.replaceAll("'",'&apos;');
    output = output.replaceAll('"','&quot;');

    return output;

}
// Debug Function for password validation 
//
// Returns true or false based on if the password is valid
//
exports.passwordCheckDebug = (password)=>{

    let debugPassword = "12345678"
    
    return password == debugPassword

}


function fragmentCompletionFunction(filepath, data){

    let test = {test:1}

    if (data == null || filepath == "")
        return ""

    let text = ""
    try {
        text = fs.readFileSync(filepath).toString()
    } catch (error) {
        console.log(0)
        return ""
    }
    console.log(4)

    //return "x"

    for (element of Object.keys(data)){

        text = text.replaceAll('${'+element+'}', data[element]);
        console.log(test[element])

    }

    return text;

} 

function fragmentStreamlined(filename, data){

    if (data == null || filename == "")
        return ""
    
    return fragmentCompletionFunction("./UI/HTML/fragments/" + filename, data)

}

exports.fragmentCompletionFunction = fragmentCompletionFunction;
exports.fragmentStreamlined = fragmentStreamlined;
