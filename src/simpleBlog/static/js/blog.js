function getFormJson(frm) {
    var o = {};
    var a = $(frm).serializeArray()
    $.each(a, function () {
        o[this.name] = this.value || '';
    })
    return o;
}

function ajaxSubmit(frm, fn) {
    var dat = getFormJson(frm)
    return $.ajax({
        url: frm.action,
        type: frm.method,
        data: dat,
        success: fn
    })
}


function ajaxx(frm, fn) {
    var dat = getFormJson(frm)
    return $.ajax({
        url: "/file/imgupload",
        contentType: "multipart/form-data",
        type: frm.method,
        success: fn
    })
}