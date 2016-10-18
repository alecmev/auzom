import { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

import Button from '../components/Button';
import Checkbox from '../components/Checkbox';
import Datetime from '../components/Datetime';
import Input from '../components/Input';
import Select from '../components/Select';
import Textarea from '../components/Textarea';

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,
}), actions)
export default class New extends Component {
  static propTypes = {
    amAdmin: PropTypes.bool,

    messagePush: PropTypes.func.isRequired,

    lighter: PropTypes.bool,
    isAdminOnly: PropTypes.bool,
    resourceName: PropTypes.string,
    fields: PropTypes.object.isRequired, // eslint-disable-line
    callForAction: PropTypes.string,
    link: PropTypes.object, // eslint-disable-line
    createResource: PropTypes.func,
    onSuccess: PropTypes.func,
    noMessage: PropTypes.bool,
    an: PropTypes.bool,
  };

  // TODO: turn fields and link into proper shapes

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentDidMount() {
    this.r[Object.keys(this.props.fields)[0]].focus();
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps);
  }

  load({ amAdmin, isAdminOnly, resourceName }) {
    if (isAdminOnly && !amAdmin) {
      this.props.messagePush(
        `only admins can create ${resourceName || 'this resource'}`, true,
      );
      this.context.router.push('/');
      return;
    }
  }

  submit = () => {
    const {
      messagePush, resourceName, fields, createResource, onSuccess, noMessage,
      an,
    } = this.props;

    const body = {};
    for (const [field, { type }] of Object.entries(fields)) {
      if (type === 'hidden') {
        body[field] = fields[field].value;
      } else {
        const value = this.r[field].state.value;
        if (value === undefined) {
          body[field] = value;
        } else if (type === 'select') {
          body[field] = value.value;
        } else if (type === 'multi') {
          body[field] = value.map(x => x.value);
        } else {
          body[field] = value;
        }
      }
    }

    createResource(body, noMessage ?
      onSuccess : (...args) => {
        messagePush(`created ${an ? 'an' : 'a'} ${resourceName || 'resource'}`);
        onSuccess && onSuccess(...args);
      }
    );
  };

  render() {
    const { resourceName, fields, callForAction, link } = this.props;

    let lastNonBool;
    let lastBools = [];
    for (const [field, { type }] of Object.entries(fields)) {
      if (type !== 'boolean') {
        lastNonBool = field;
        lastBools = [];
      } else lastBools.push(field);
    }

    const form = [];
    const defaultDate = (new Date()).setMinutes(0, 0, 0);
    for (const [
      field, { type, name, value, options },
    ] of Object.entries(fields)) {
      if (type === 'hidden') {
        // continue
      } else if (type === 'date') {
        form.push(
          <Datetime
            key={field}
            ref={this.rcb(field)}
            label={`${name}...`}
            justDate
            value={value === 'now' ? defaultDate : null}
          />,
        );
      } else if (type === 'datetime') {
        form.push(
          <Datetime
            key={field}
            ref={this.rcb(field)}
            label={`${name}...`}
            value={value === 'now' ? defaultDate : null}
          />,
        );
      } else if (type === 'select') {
        form.push(
          <Select
            key={field}
            ref={this.rcb(field)}
            options={options}
            label={`${name}...`}
          />,
        );
      } else if (type === 'multi') {
        form.push(
          <Select
            key={field}
            ref={this.rcb(field)}
            options={options}
            label={`${name}...`}
            multi
          />,
        );
      } else if (type === 'boolean') {
        form.push(
          <Checkbox
            key={field}
            ref={this.rcb(field)}
            label={name}
            value={value}
            onEnter={lastBools.includes(field) ? this.submit : undefined}
          />,
        );
      } else if (type === 'textarea') {
        form.push(
          <Textarea
            key={field}
            ref={this.rcb(field)}
            label={name}
            value={value}
          />,
        );
      } else {
        form.push(
          <Input
            key={field}
            ref={this.rcb(field)}
            type={['email', 'password'].includes(type) ? type : 'text'}
            label={name}
            value={value}
            onEnter={field === lastNonBool ? this.submit : undefined}
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
            text={callForAction || `create new ${resourceName || 'resource'}`}
            onClick={this.submit}
          />
          {link &&
            <Link
              className={this.cn({ d: 'link' })}
              to={link.to}
            >
              {link.name}
            </Link>
          }
        </div>
      </div>
    );
  }
}
