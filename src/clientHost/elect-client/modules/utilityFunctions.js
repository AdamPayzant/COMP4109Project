// Function that returns text after the html escape characrers have been replaced with html codes
//
// List of all escape characters (in case we want to add more): https://www.werockyourweb.com/html-escape-characters/
//
exports.sanatizeText = (input)=>{

    let output = input.toString();

    if(!output){
        return ""
    }

    output = output.replaceAll('<','&lt;');
    output = output.replaceAll('>','&gt;');
    output = output.replaceAll("'",'&apos;');
    output = output.replaceAll('"','&quot;');

    return output;

}