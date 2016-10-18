import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class PasswordReset extends Component {
  static propTypes = {
    createOTP: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleSuccess = () => {
    browserHistory.push('/');
    this.props.messagePush('password reset link sent, check your email');
  };

  render() {
    return (
      <New
        fields={{
          email: { type: 'email', name: 'Your email' },
        }}
        callForAction="SEND PASSWORD RESET LINK"
        createResource={this.props.createOTP}
        onSuccess={this.handleSuccess}
        noMessage
      />
    );
  }
}
