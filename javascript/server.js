

var queryTX = require("./queryTx.js");
var queryTU = require("./queryUsr.js");
var invoke = require("./invokeExport.js");
var express = require("express");
var app = express();
var fs = require("fs");
var bodyParser = require("body-parser")
var swig=require('swig');
const { request } = require("http");

// add by cwh
app.use(express.static(__dirname + '/html'));
app.set('views', './');
//设置html模板渲染引擎
app.engine('html', swig.renderFile);
//设置渲染引擎为html
app.set('view engine','html');

app.use(bodyParser.urlencoded({
    extended: true
}));

// add by cwh
app.get('/', function (request, response) {
  response.redirect('/main');
});

app.get('/main', function (request, response) {
    response.render('html/main', {result_dic_list: []});
});

app.get('/deploy', function (request, response) {
  response.render('html/deploy', {});
});

app.post('/main', function (request, response) {
    console.log(request.body);
    func = request.body.func;
    radio_type = request.body.radio_type;
    check_settlement = request.body.check_settlement;
    query = request.body.query;
    if (func == 'deploy'){
      console.log("deploy");
      console.log(request.body);
      tx = request.body.tx;
      usr = request.body.usr;
      caas = request.body.caas;
      startTime = request.body.startTime;
      endTime = request.body.endTime;
      signals = ""
      VcState = "1"
      ID_CO = tx + "," + usr + "," + caas + "," + startTime + "," + endTime
      console.log(ID_CO);
      invoke.invokecc(func, [ID_CO]).then((result) => {
        result_dic = {
          "Tx": tx,
          "Usr": usr,
          "Caas": caas,
          "StartTime": startTime,
          "EndTime": endTime,
          "VcState": VcState
        }
        result_dic_list = [];
        console.log(result);
       
        result_dic_list.push(result_dic);
        tmp_signal = tx + "," + "1" + "," + "100" + "," + "0.5"
        invoke.invokecc("control", [tmp_signal]).then((tmp_signal_result) => {
          console.log(tmp_signal_result);
        });
        console.log(result_dic_list);
        response.render('html/main', {
          result_dic_list: result_dic_list
        });
      });

    } else if (func == 'control') {
      tx = request.query.tx;
      t_i = request.query.t_i;
      p_i = request.query.p_i;
      soc_i = request.query.soc_i;
      ID_CO = tx + "," + t_i + "," + p_i + "," + soc_i
      console.log(ID_CO);
      invoke.invokecc(func, [ID_CO]).then((result) => {
        console.log(result);
      });
    } else if (radio_type == 0) { // usr
        ID_CO = query;
        console.log(ID_CO);
        queryTU.queryCO(ID_CO).then((result) => {
            result_list = result.split("}{");
            result_dic_list = [];
            result_list.forEach(function(item, index){
            if (index == 0){
              if (result_list.length > 1){
                item = item + "}";
              }
            } else if (index == result_list.length - 1){
                item = "{" + item;
            } else {
                item = "{" + item + "}";
            };
            result_dic_list.push(JSON.parse(item));
            });
            console.log(result_dic_list);
            response.render('html/main',{
            result_dic_list: result_dic_list
            });
        });
    } else if (radio_type == 1) { 
        if (check_settlement == 'on') { // settlement
            ID_CO = query;
            console.log(ID_CO);
            func = "settlement";
            invoke.invokecc(func, [ID_CO]).then((result) => {
              console.log(result);
              settlement_ok = False
              if (result == 200){
                settlement_ok = True
              }
              queryTX.queryCO(ID_CO).then((result) => {
                result_dic_list = [];
              if (result.length > 2){
              result_dic = JSON.parse(result);
              result_dic_list.push(result_dic);
              } 
              
              console.log(result_dic_list);
              response.render('html/main',{
                  result_dic_list: result_dic_list,
                  settlement_ok: settlement_ok
                  });
              });
            });
            
        } else { // tx
          console.log("tx");

            ID_CO = query;
            console.log(ID_CO);
            queryTX.queryCO(ID_CO).then((result) => {
              result_dic_list = [];
            if (result.length > 2){
            result_dic = JSON.parse(result);
            result_dic_list.push(result_dic);
            } 
            
            console.log(result_dic_list);
            response.render('html/main',{
                result_dic_list: result_dic_list
                });
            });
        };
    };
});

app.listen(1025);
