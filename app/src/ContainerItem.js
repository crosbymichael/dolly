import React, { Component } from 'react';
import SmoothieGraph from './SmoothieGraph';
import classnames from 'classnames';

export default class ContainerItem extends Component {
  render() {
    var overrideButtonStyle = {
      color: 'white',
      backgroundColor: '#3c5164'
    };
    var barStyle = {
      transitionDuration: '300ms',
      width: this.props.fillPct + '%'
    };
    var barClassnames = classnames({
      'red': this.props.fillPct > 90,
      'orange': this.props.fillPct <= 90 && this.props.fillPct > 70,
      'yellow': this.props.fillPct <= 70 && this.props.fillPct > 50,
      'olive': this.props.fillPct <= 50,
      'ui': true,
      'progress': true
    });
    return (
      <div className="card">
        <div className="content">
          <div className="header">
            {this.props.name || 'unknown'}
            <div className="right">
              <div className={barClassnames} data-percent={this.props.fillPct}>
                <div className="bar" style={barStyle}></div>
              </div>
            </div>
          </div>
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
  name: React.PropTypes.string,
  fillPct: React.PropTypes.number
}
