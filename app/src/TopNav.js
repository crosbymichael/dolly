import React, { Component } from 'react';

export default class TopNav extends Component {
  render() {
    return (
      <div className="ui pointing menu">
        <div className="right menu">
          <div className="item">
            <div className="statistic">
              <div className="value">
                2,204
              </div>
              <div className="label">
                Requests
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
