import React, { Component } from 'react';

export default class TopNav extends Component {
  constructor(props) {
    super(props);
    this.state = {
      totalRequests: 0
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
      totalRequests: result.totalRequests
    });
  }
  render() {
    return (
      <div className="ui pointing menu">
        <div className="right menu">
          <div className="item">
            <div className="statistic">
              <div className="value">
                {this.state.totalRequests}
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
