import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

@connect(createStructuredSelector({
  canGoBack: selectors.canGoBack,
}), actions)
export default class Login extends Component {
  static propTypes = {
    canGoBack: PropTypes.bool.isRequired,
    createSession: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleSuccess = () => {
    if (this.props.canGoBack) browserHistory.goBack();
    else browserHistory.push('/');
    this.props.messagePush('logged in');
  };

  render() {
    return (
      <New
        fields={{
          email: { type: 'email', name: 'Email' },
          password: { type: 'password', name: 'Password' },
          remember: { type: 'boolean', name: 'Stay logged in', value: true },
        }}
        callForAction="log in"
        link={{ to: '/password-reset', name: 'Forgot your password?' }}
        createResource={this.props.createSession}
        onSuccess={this.handleSuccess}
        noMessage
      />
    );
  }
}
