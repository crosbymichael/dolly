import React, { Component } from 'react';
import SmoothieGraph from './SmoothieGraph';

export default class ContainerItem extends Component {
  render() {
    var overrideButtonStyle = {
      color: 'white',
      backgroundColor: '#3c5164'
    };
    return (
      <div className="card">
        <div className="content">
          <div className="header">{this.props.name || 'unknown'}</div>
          <div className="description">
            <SmoothieGraph />
          </div>
        </div>
        <div className="extra content">
          <div className="ui two buttons">
            <div className="ui basic green button">Start</div>
            <div className="ui basic blue button">Clone</div>
          </div>
      </div>
    </div>
    );
  }
}

ContainerItem.propTypes = {
  name: React.PropTypes.string
}
