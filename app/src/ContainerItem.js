import React, { Component } from 'react';
import LineGraph from './LineGraph';

export default class ContainerItem extends Component {
  render() {
    return (
      <div className="card">
        <div className="content">
          <div className="header">modest_thompson</div>
          <div className="description">
            <LineGraph />
          </div>
        </div>
        <div className="ui bottom attached button">
          <i className="play icon"></i>
          Start
        </div>
      </div>
    );
  }
}
