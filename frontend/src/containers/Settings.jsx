import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

import NotFound from './NotFound';
import Loading from './Loading';

import Button from '../components/Button';
import Checkbox from '../components/Checkbox';
import Datetime from '../components/Datetime';
import Input from '../components/Input';
import Textarea from '../components/Textarea';

@connect(createStructuredSelector({
  isLoading: selectors.isLoading,
  amAdmin: selectors.amAdmin,
}), actions)
export default class Settings extends Component {
  static propTypes = {
    isLoading: PropTypes.bool.isRequired,
    amAdmin: PropTypes.bool,

    messagePush: PropTypes.func.isRequired,

    lighter: PropTypes.bool,
    isAdminOnly: PropTypes.bool,
    resourceName: PropTypes.string,
    resource: ImmutablePropTypes.map,
    fields: PropTypes.object.isRequired, // eslint-disable-line
    callForAction: PropTypes.string,
    updateResource: PropTypes.func.isRequired,
    onSuccess: PropTypes.func,
  };

  // TODO: turn fields into a proper shape

  state = {
    resource: null,
    values: null,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  getModified() {
    return Object.keys(this.props.fields).reduce((r, field) => ({
      ...r, [field]: this.state.values[field] !== this.state.resource[field],
    }), {});
  }

  load({ amAdmin, isAdminOnly, resourceName, resource, fields }, prevProps) {
    if (isAdminOnly && !amAdmin) {
      this.props.messagePush(
        `only admins can change ${resourceName} settings`, true,
      );
      browserHistory.push('/');
      return;
    }

    // no need in resetting the state, it won't be visible anyway
    if (!resource) return;

    if (!prevProps || resource !== prevProps.resource) {
      const resourceJS = resource.toJS();
      for (const [field, { type }] of Object.entries(fields)) {
        if (type === 'date' || type === 'datetime') {
          if (resourceJS[field] !== null) {
            // Normalizing the date/time format, as to be able to detect
            // modifications properly. One important nitpick is that we lose
            // nanoseconds permanently here, thanks to JS...
            resourceJS[field] = moment(resourceJS[field]).toISOString();
          }
        } else if (type === 'number') {
          if (typeof resourceJS[field] === 'number') {
            resourceJS[field] = resourceJS[field].toString();
          }
        }
      }

      this.setState({
        resource: resourceJS,
        // would be a shame to lose all modifications due to some late update...
        values: this.state.values || resourceJS,
      });
    }
  }

  submit = () => {
    const { messagePush, resourceName, updateResource, onSuccess } = this.props;
    const { values } = this.state;

    const modified = this.getModified();
    const body = {};
    let modifiedAtAll = false;
    for (const [field, isModified] of Object.entries(modified)) {
      if (isModified) {
        body[field] = values[field];
        modifiedAtAll = true;
      }
    }

    if (!modifiedAtAll) return;
    updateResource(
      values.id, body,
      onSuccess ? onSuccess.bind(null, modified) :
        () => messagePush(`${resourceName || 'resource'} settings saved`),
    );
  };

  handleChange = (x, field) => this.setState({
    values: { ...this.state.values, [field]: x },
  });

  render() {
    const {
      isLoading, resourceName, resource, fields, callForAction,
    } = this.props;
    const { values } = this.state;

    if (!resource) {
      return isLoading ?
        <div className={this.cni()}><Loading /></div> :
        <NotFound />;
    }

    const modified = this.getModified();
    let modifiedAtAll = false;
    const form = [];
    for (const [field, isModified] of Object.entries(modified)) {
      if (isModified) modifiedAtAll = true;
      const { type, name } = fields[field];
      if (type === 'date') {
        form.push(
          <Datetime
            key={field}
            justDate
            value={values[field]}
            onChange={this.handleChange}
            parentField={field}
            label={name}
            modified={isModified}
          />,
        );
      } else if (type === 'datetime') {
        form.push(
          <Datetime
            key={field}
            value={values[field]}
            onChange={this.handleChange}
            parentField={field}
            label={name}
            modified={isModified}
          />,
        );
      } else if (type === 'boolean') {
        form.push(
          <Checkbox
            className={this.cn({ d: 'checkbox', s: isModified && 'modified' })}
            key={field}
            label={name}
            value={values[field]}
            onChange={this.handleChange}
            parentField={field}
            modified={isModified}
          />,
        );
      } else if (type === 'textarea') {
        form.push(
          <Textarea
            key={field}
            value={values[field]}
            onChange={this.handleChange}
            parentField={field}
            label={name}
            modified={isModified}
          />,
        );
      } else {
        form.push(
          <Input
            key={field}
            type={['email', 'password'].includes(type) ? type : 'text'}
            value={values[field]}
            onChange={this.handleChange}
            parentField={field}
            label={name}
            modified={isModified}
          />,
        );
      }
    }

    return (
      <div className={this.cni({ m: this.props.lighter && 'lighter' })}>
        <div className={this.cn({ d: 'inner' })}>
          {form}
          <Button
            block
            className={this.cn({ d: 'submit' })}
            type="important"
            text={
              callForAction || `save ${resourceName || 'resource'} settings`
            }
            onClick={this.submit}
            disabled={!modifiedAtAll}
          />
        </div>
      </div>
    );
  }
}
