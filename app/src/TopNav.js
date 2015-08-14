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
    var overrideNavStyle = {
      color: 'white',
      backgroundColor: '#3c5164',
      borderRadius: '0'
    };
    var overrideStatStyle = {
      color: 'white'
    };
    var menuStyle = {
      borderLeft: '1px white solid'
    };
    var iconStyle = {
      marginTop: '0.25rem',
      marginLeft: '0.75rem'
    };
    var headingStyle = {
      margin: '1rem',
      color: 'white'
    };
    return (
      <div className="ui pointing menu" style={overrideNavStyle}>
        <div>
          <img src="styles/dist/images/mini-logo.svg" style={iconStyle}/>
        </div>
        <div>
          <a href='https://dolly.dockerproject.org'>
            <h3 style={headingStyle}>dolly.dockerproject.org</h3>
          </a>
        </div>
      </div>
    );
  }
}
