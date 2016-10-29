(function () {
  $(window).load(function() {
    //console.log("window loaded");
    updateRequestHistory();
  });
  // $(document).ready(function() {
  //   console.log("document loaded");
  //   connectWebsocket();
  // });
})();

function getMessage(topic, partition, offset) {
  console.log("getMessage => ", topic, partition, offset)
  //alert(topic + '|' + partition + '|' + offset)
  var urlRequestDetails = "/api/request/payload/"+topic+"/"+partition+"/"+offset;
  $.ajax({
     url: urlRequestDetails,
     contentType: "application/octet-stream; charset=utf-8",
     type: 'GET',
     success: function(data) {
        $("#titleModal").text("Payload Message")
        $("#messageModal").text(data);
        $("#myModal").modal();
     },
     error: function(data) {
         console.log("error: ", data);
     }
 });
}

function getError(err) {
  $("#titleModal").text("Error Message");
  $("#messageModal").text(err);
  $("#myModal").modal();
}

function tableDataAddItem(item) {
  var status_green = "#88CF85";
  var status_red = "#F87070";
  var duration_green = "#48B444";
  var duration_orange = "#F4B27B";
  var duration_red = "#F53636";

  var status_color = item.status == 200 ? status_green : status_red;
  var status_class = item.status == 200 ? "success" : "fail";
  var status_style = "";
  if ((statusFilter == "success" && item.status != 200) || (statusFilter == "fail" && item.status == 200)) {
    status_style = "style=\"display: none;\""
  }
  var duration_color;
  if (item.duration_ms < 1000) {
    duration_color = duration_green;
  } else if (item.duration_ms >= 1000 && item.duration_ms < 2000) {
    duration_color = duration_orange;
  } else {
    duration_color = duration_red;
  }
  $("#tableDataHeader").after(
    "<tbody id=\""+item.rid+"\" class=\"bodycontent "+status_class+"\" "+status_style+">"+
      "<tr class=\"trcontent\">"+
        "<th class=\"text-center\">"+
          "<a href=\"javascript:void(0)\"><span id=\"icon-"+item.rid+"\" class=\"glyphicon glyphicon-plus\" aria-hidden=\"true\" onclick=\"getRequestDetails('"+item.rid+"')\"></span></a>"+
        "</th>"+
        "<td>"+item.time+"</td>"+
        "<td>"+item.rid+"</td>"+
        "<td>"+item.supervisor+"</td>"+
        "<td>"+item.module+"</td>"+
        "<td><span style=\"color:"+duration_color+"\">"+parseTime(item.duration_ms)+"</span></td>"+
        "<td bgcolor=\""+status_color+"\">"+item.status+"</td>"+
      "</tr>"+
    "</tbody>");
    var count = 0;
    var trs_success = $("#tableData").find(".bodycontent.success");
    var trs_fail = $("#tableData").find(".bodycontent.fail");
    if (statusFilter == "success") {
      count = trs_success.length;
    } else if (statusFilter == "fail") {
      count = trs_fail.length;
    } else {
      count = trs_success.length + trs_fail.length;
    }
    tableDataAppendCount(count);
}

function tableDataAddInfo(msg) {
  $("#tableDataHeader").after(
    "<tbody class=\"bodycontent\">"+
    "  <tr class=\"trcontent\">"+
    "    <td colspan=\"8\" class=\"text-center\"><h6>"+msg+"</h6></td>"+
    "  </tr>"+
    "</tbody>");
}

function tableDataAppendCount(count) {
  $("#tableDataCount").remove();
  $("#tableData").append(
    "<tbody id=\"tableDataCount\" class=\"bodycontent\">"+
    "  <tr class=\"trcontent\">"+
    "    <td colspan=\"8\" class=\"text-center\"><h6>Found "+count+" records</h6></td>"+
    "  </tr>"+
    "</tbody>");
}

function updateRequestHistory() {
  var blocksize = $("#datasize").val()
  var urlRequestHistory = "/api/request/history?blocksize="+blocksize;
  $.getJSON( urlRequestHistory, {
    format: "json"
  })
  .done(function(data) {
    var count = data.length
    //var data_reverse = data.slice(0).reverse();
    $(".bodycontent").remove()
    $("#statusDropDown").text("Status: all")
    $.each(data, function(i, item) {
      tableDataAddItem(item);
    });
    tableDataAppendCount(count);
    $('#reqHistButton').focus();
  });
}
//from miliseconds to seconds
function parseTime(duration_ms){
  var time;

  if (duration_ms < 1000) {
    time = Number(parseFloat(duration_ms)).toFixed(0)+"ms"
  } else if (duration_ms >= 1000 && duration_ms < 2000) {
    time = Number(parseFloat(duration_ms)/1000).toFixed(2)+"s"
  } else {
    var sec = parseFloat(duration_ms)/1000;
    var min = Math.floor(sec/60)
    var sec = Math.floor(sec % 60)
    time = min+"m "+sec+"s"
  }
  return time
}

function getRequestDetails(rid) {

  if ($("#tbase-"+rid).length) {
    $("#icon-"+rid).removeClass("glyphicon-minus").addClass("glyphicon-plus");
    $("#tbase-"+rid).remove();
    return;
  }

  var urlRequestDetails = "/api/request/details/"+rid;
  console.log("getRequestDetails: ", urlRequestDetails)
  $.getJSON(urlRequestDetails, {
    format: "json"
  })
  .done(function(data) {
    var details =
      "<table class=\"table table-condensed table-bordered table-striped table-hover\">"+
      "    <thead>"+
      "      <tr>"+
      "        <th>Step</th>"+
      "        <th>RequestId</th>"+
      "        <th>Host</th>"+
      "        <th>Service</th>"+
      "        <th>Message</th>"+
      "        <th>Topic</th>"+
      "        <th>Partition</th>"+
      "        <th>Offset</th>"+
      "        <th>Msg</th>"+
      "        <th>Error</th>"+
      "        <th>Duration [ms]</th>"+
      "        <th>Status</th>"+
      "        </tr>"+
      "    </thead>";

    // get message status
    var globalMessageStatus = 0;
    $.each(data, function(i, item) {
      if (item.status != 0 && globalMessageStatus == 0) {
        globalMessageStatus = item.status;
      }
    });
    var status_class = globalMessageStatus == 200 ? "details_success" : "details_fail";
    $.each(data, function(i, item) {
      var msgHidden = 'hidden'
      var errHidden = 'hidden'
      if (item.topic.length > 0 && (item.partition > 0 || item.offset > 0))
        msgHidden = ''
      if (item.status != 200 && item.status != 0 && item.message_error.length > 0)
        errHidden = ''

      // optional messages
      if (item.step === 100) {
        details +=
          "<tbody id=\"tb-"+item.rid+"\" class=\"bodycontent "+status_class+"\">"+
          "  <tr class=\"trcontent\">"+
          "    <td colspan=\"12\" class=\"text-center\">additional messages</td>"+
          "  </tr>"+
          "</tbody>";
      }
      var progress = item.progress || item.step;
      details +=
        "<tbody id=\"tb-"+item.rid+"\" class=\"bodycontent"+status_class+"\">"+
        "  <tr class=\"trcontent\">"+
        "    <td>"+progress+"</td>"+
        "    <td>"+item.rid+"</td>"+
        "    <td>"+item.host+"</td>"+
        "    <td>"+item.supervisor+item.module+"</td>"+
        "    <td>"+item.message+"</td>"+
        "    <td>"+item.topic+"</td>"+
        "    <td>"+item.partition+"</td>"+
        "    <td>"+item.offset+"</td>"+
        "    <td><button type=\"button\" class=\""+msgHidden+" btn btn-default btn-sm\" onclick=\"getMessage('"+item.topic+"',"+item.partition+","+item.offset+")\">"+
        "        <span class=\"glyphicon glyphicon-file\"></span></button>"+
        "    </td>"+
        "    <td><button type=\"button\" class=\""+errHidden+" btn btn-default btn-sm\" onclick=\"getError('"+item.message_error+"')\">"+
        "        <span class=\"glyphicon glyphicon-exclamation-sign\"></span></button>"+
        "    </td>"+
        "    <td>"+parseTime(item.duration_ms)+"</td>"+
        "    <td>"+item.status+"</td>"+
        "  </tr>"+
        "</tbody>";
    });
    details += "</table>";
    $("#icon-"+rid).removeClass("glyphicon-plus").addClass("glyphicon-minus");
    $("#"+rid).append(
      "    <tr id=\"tbase-"+rid+"\" class=\"trcontent\">"+
      "        <th class=\"text-center\">"+
      "        </th>"+
      "        <td colspan=\"8\">"+details+"</td>"+
      "    </tr>")
  });
}

var websocket;
function connectWebsocket() {

  function resetWs() {
    websocket.onmessage = function () {};
    websocket.onclose = function () {};
    websocket.onopen = function () {};
    websocket.close();
    websocket = null;
  }

  if (websocket != null) {
    resetWs();
    tableDataAddInfo("Websocket connection closed");
    $("#wsButton").removeClass('active').removeClass('btn-primary');
    return;
  }

  if (window["WebSocket"]) {
      var host = window.location.host;
      websocket = new WebSocket("ws://" + host + "/ws");
      websocket.onopen = function(evt) {
        tableDataAddInfo("Websocket connection opened");
        $("#wsButton").addClass('active').addClass('btn-primary');
      }
      websocket.onclose = function(evt) {
        tableDataAddInfo("Websocket connection closed");
        $("#wsButton").removeClass('active').removeClass('btn-primary');
        resetWs();
      };
      websocket.onmessage = function(evt) {
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
          var item = JSON.parse(messages[i]);
          tableDataAddItem(item);
        }
      };
  } else {
    tableDataAddInfo("Your browser does not support WebSockets");
  }
}

var statusFilter = null;
function filterTableDataByClass(filterClass) {
  if (statusFilter == filterClass) {
    console.log("same filter: ", statusFilter);
    return;
  }
  statusFilter = filterClass;

  var count = 0;
  var trs_all = $("#tableData").find(".bodycontent");
  count = trs_all.length - 1;
  trs_all.hide();
  if (filterClass != null) {
    var trs = $("#tableData").find(".bodycontent."+filterClass);
    count = trs.length;
    var trs_children = $("#tableData").find(".bodycontent.details_"+filterClass);
    trs_children.show();
    trs.show();
    $("#statusDropDown").text("Status: "+filterClass)
  } else {
    trs_all.show();
    $("#statusDropDown").text("Status: all")
  }
  tableDataAppendCount(count);
}
