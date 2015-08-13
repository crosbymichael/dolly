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
    var result = {
       "totalRequests" : 2020,
       "servers" : [
          {
             "statsEndpoint" : "http://127.0.0.3001",
             "fill" : 66.6,
             "name" : "linuxcon1"
          },
          {
             "fill" : 100,
             "name" : "linuxcon1",
             "statsEndpoint" : "http://127.0.0.3002"
          }
       ]
    }

    this.setState({
      totalRequests: result.totalRequests,
      containersList: result.servers
    });
  }
  renderContainers() {
    var _makeServerItem = function(server) {
      return <ContainerItem name={server.name} />;
    }
    return this.state.containersList.map(_makeServerItem);
  }
  render() {
    return (
      <div className="ui cards six column padded grid">
        {this.renderContainers()}
      </div>
    );
  }
}
