var deploy = require("./deploy.js"); //queryCO
var queryTX = require("./queryTx.js");
var queryTU = require("./queryUsr.js");
var control = require("./control.js");
var settlement = require("./settlement.js");
var express = require("express");
var app = express();
var fs = require("fs");
var bodyParser = require("body-parser")
var swig=require('swig');

app.set('views', './');
//设置html模板渲染引擎
app.engine('html', swig.renderFile);
//设置渲染引擎为html
app.set('view engine','html');

app.use(bodyParser.urlencoded({
    extended: true
}));


app.get('/main', function (request, response) {
    response.render('main', {});
});

app.get('/deploy', function (request, response) {
  response.render('main', {});
});

app.post('/queryTx', function (request, response) {
    ID_CO = request.body.tx;
    console.log(ID_CO);
    queryTX.queryCO(ID_CO).then((result) => {
      result_list = result.split(",");
      response.render('main',{
        tx: result_list[0],
        usr: result_list[1],
        caas: result_list[2],
        startTime: result_list[3],
        endTime: result_list[4],
        signals: result_list[5],
        pay: result_list[6]
      });
    });
});

app.post('/queryUsr', function (request, response) {
  ID_CO = request.body.usr;
  console.log(ID_CO);
  queryTU.queryCO(ID_CO).then((result) => {
    result_list = result.split(",");
    response.render('main',{
      tx: result_list[0],
      usr: result_list[1],
      caas: result_list[2],
      startTime: result_list[3],
      endTime: result_list[4],
      signals: result_list[5],
      pay: result_list[6]
    });
  });
});

app.post('/deploy', function (request, response) {
  tx = request.body.tx;
  usr = request.body.usr;
  caas = request.body.caas;
  startTime = request.body.startTime;
  endTime = request.body.endTime;
  ID_CO = tx + "," + usr + "," + caas + "," + startTime + "," + endTime
  console.log(ID_CO);
  deploy.queryCO(ID_CO).then((result) => {
    result_list = result.split(",");
    response.render('main',{
      tx: result_list[0],
      usr: result_list[1],
      caas: result_list[2],
      startTime: result_list[3],
      endTime: result_list[4],
      signals: result_list[5],
      pay: result_list[6]
    });
  });
});

app.post('/control', function (request, response) {
  tx = request.body.tx;
  t_i = request.body.t_i;
  p_i = request.body.p_i;
  soc_i = request.body.soc_i;
  ID_CO = tx + "," + t_i + "," + p_i + "," + soc_i
  console.log(ID_CO);
  control.queryCO(ID_CO).then((result) => {
    result_list = result.split(",");
    response.render('main', {
      tx: result_list[0],
      usr: result_list[1],
      caas: result_list[2],
      startTime: result_list[3],
      endTime: result_list[4],
      signals: result_list[5],
      pay: result_list[6]
    });
  });
});

app.post('/settlement', function (request, response) {
  ID_CO = request.body.tx;
  console.log(ID_CO);
  settlement.queryCO(ID_CO).then((result) => {
    result_list = result.split(",");
    response.render('main',{
      tx: result_list[0],
      usr: result_list[1],
      caas: result_list[2],
      startTime: result_list[3],
      endTime: result_list[4],
      signals: result_list[5],
      pay: result_list[6]
    });
  });
});

app.listen(8080);
