import React, { Component } from 'react';
import SmoothieGraph from './SmoothieGraph';
import classnames from 'classnames';
import numeral from 'numeral';
import request from 'superagent';

export default class ContainerItem extends Component {
  handleStart(evt) {
    console.log(this);
    console.log(evt);
    evt.preventDefault();
    request
      .post('http://localhost:8765/' + this.props.name + '/start')
      .send({})
      .end(function(err, res){
        // Calling the end function will send the request
        if (err) {
          console.log(err);
        }
      });
  }
  handleClone(evt) {
    evt.preventDefault();
    request
      .post('http://localhost:8765/' + this.props.name + '/clone')
      .send({})
      .end(function(err, res){
        // Calling the end function will send the request
        if (err) {
          console.log(err);
        }
      });
  }
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
      'green': this.props.fillPct > 90,
      'yellow': this.props.fillPct <= 90 && this.props.fillPct > 70,
      'orange': this.props.fillPct <= 70 && this.props.fillPct > 50,
      'red': this.props.fillPct <= 50,
      'ui': true,
      'progress': true
    });
    return (
      <div className="card four wide column">
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
            <div className="ui horizontal statistic">
              <div className="value">
                {numeral(this.props.responseTime).format('0.00')}
              </div>
              <div className="label">
                rps
              </div>
            </div>
          </div>
        </div>
        <div className="extra content">
          <div className="ui two buttons">
            <div className="ui basic green button" onClick={this.handleStart.bind(this)}>Start</div>
            <div className="ui basic blue button" onClick={this.handleClone.bind(this)}>Clone</div>
          </div>
      </div>
    </div>
    );
  }
}

ContainerItem.propTypes = {
  name: React.PropTypes.string,
  fillPct: React.PropTypes.number,
  responseTime: React.PropTypes.number
}
