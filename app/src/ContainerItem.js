import React, { Component } from 'react';
import SmoothieGraph from './SmoothieGraph';

export default class ContainerItem extends Component {
  render() {
    return (
      <div className="card">
        <div className="content">
          <div className="header">{this.props.name || 'unknown'}</div>
          <div className="description">
            <SmoothieGraph />
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

ContainerItem.propTypes = {
  name: React.PropTypes.string
}
