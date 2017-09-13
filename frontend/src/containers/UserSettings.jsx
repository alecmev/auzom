import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import Settings from './Settings';

@connect(createStructuredSelector({
  myId: selectors.myId,
  me: selectors.me,
}), actions)
export default class UserSettings extends Component {
  static propTypes = {
    myId: PropTypes.string,
    me: ImmutablePropTypes.map,

    updateUser: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps);
  }

  load({ myId }) {
    if (!myId) {
      this.props.messagePush(
        'you need to be logged in to change your settings', true,
      );
      browserHistory.push('/login');
    }
  }

  handleSuccess = m => this.props.messagePush(
    `settings saved${m.email ? ', check your email' : ''}`,
  );

  render() {
    return (
      <Settings
        resourceName="user"
        resource={this.props.me && this.props.me.set('password', '')}
        fields={{
          email: { type: 'email', name: 'Email' },
          password: { type: 'password', name: 'New password' },
          nickname: { name: 'Nickname' },
          fullname: { name: 'Full name' },
          gravatarEmail: { type: 'email', name: 'Gravatar email' },
        }}
        callForAction="save settings"
        updateResource={this.props.updateUser}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
