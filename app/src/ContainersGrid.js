import React, { Component } from 'react';
import ContainerItem from './ContainerItem';
import request from 'superagent';

export default class ContainersGrid extends Component {
  constructor(props) {
    super(props);
    this.state = {
      totalRequests: 0,
      containersList: []
    };
  }
  componentDidMount() {
    var result;
    var _this = this;

    (function poll(){
       setTimeout(function(){
         request.get('http://localhost:8765/')
           .end(function(err, res){
               if (err) {
                 console.log(err);
               } else {
                 result = JSON.parse(res.text);
                 console.log(result);
                 _this.setState({
                   totalRequests: result.totalRequests,
                   containersList: result.servers
                 });
                 poll();
               }
           })}, 2000);
    })();
  }
  renderContainers() {
    var _makeServerItem = function(server) {
      return <ContainerItem name={server.name} fillPct={server.fill} responseTime={server.responseTime} />;
    }
    return this.state.containersList.map(_makeServerItem);
  }
  render() {
    return (
      <div className="ui cards padded grid">
        {this.renderContainers()}
      </div>
    );
  }
}
