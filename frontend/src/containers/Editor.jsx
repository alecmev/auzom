import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';
import Component from '../utils/Component';

import Button from '../components/Button';
import Checkbox from '../components/Checkbox';
import Datetime from '../components/Datetime';
import Input from '../components/Input';
import Img from '../components/Img';
import Markdown from '../components/Markdown';
import Textarea from '../components/Textarea';

const defValues = {
  date: null,
  datetime: null,
  boolean: false,
};

@connect(undefined, actions)
export default class Editor extends Component {
  static propTypes = {
    isNew: PropTypes.bool,
    resourceName: PropTypes.string,
    resource: ImmutablePropTypes.map,
    fields: PropTypes.object.isRequired, // eslint-disable-line
    callForAction: PropTypes.string,
    createResource: PropTypes.func,
    updateResource: PropTypes.func,
    onSuccess: PropTypes.func,
    noMessage: PropTypes.bool,

    messagePush: PropTypes.func.isRequired,
  };

  // TODO: fields -> shape

  constructor(props) {
    super(props);
    this.state = {
      original: null,
      values: !props.isNew ? null : Object.entries(props.fields).reduce(
        (r, [k, x]) => {
          let v = '';
          if (x.value !== undefined) v = x.value;
          else if (defValues[x.type] !== undefined) v = defValues[x.type];
          return { ...r, [k]: v };
        }, {},
      ),
    };
  }

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  getModified() {
    return Object.keys(this.props.fields).reduce((r, field) => ({
      ...r,
      [field]: this.props.isNew ? false :
        this.state.values[field] !== this.state.original[field],
    }), {});
  }

  load({ isNew, resource, fields }, prevProps) {
    if (isNew || !resource) return;
    if (!prevProps || resource !== prevProps.resource) {
      const resourceJS = resource.toJS();
      Object.entries(fields).forEach(([field, { type }]) => {
        if (type === 'date' || type === 'datetime') {
          if (resourceJS[field] !== null) {
            resourceJS[field] = moment(resourceJS[field]).toISOString();
          }
        } else if (type === 'number') {
          if (typeof resourceJS[field] === 'number') {
            resourceJS[field] = resourceJS[field].toString();
          }
        }
      });

      this.setState({
        original: resourceJS,
        values: this.state.values || resourceJS,
      });
    }
  }

  handleChange = (x, field) => this.setState({
    values: { ...this.state.values, [field]: x },
  });

  submit = () => {
    const {
      resourceName, createResource, updateResource, onSuccess, noMessage,
      messagePush,
    } = this.props;
    const { values } = this.state;
    if (this.props.isNew) {
      createResource(values, noMessage ?
        onSuccess : (...args) => {
          messagePush(`created a ${resourceName || 'resource'}`);
          onSuccess && onSuccess(...args);
        },
      );
    } else {
      const modified = this.getModified();
      const body = {};
      let modifiedAtAll = false;
      Object.entries(modified).forEach(([field, isModified]) => {
        if (isModified) {
          body[field] = values[field];
          modifiedAtAll = true;
        }
      });

      if (!modifiedAtAll) return;
      updateResource(values.id, body, noMessage ?
        onSuccess && onSuccess.bind(null, modified) : (...args) => {
          messagePush(`${resourceName || 'resource'} saved`);
          onSuccess && onSuccess(modified, ...args);
        },
      );
    }
  };

  render() {
    const {
      isNew, resourceName, resource, fields, callForAction,
    } = this.props;
    if (!isNew && !resource) return null;
    const { values } = this.state;

    const modified = this.getModified();
    let modifiedAtAll = false;
    const form = [];
    let left = [];
    let right = [];
    Object.entries(modified).forEach(([field, isModified]) => {
      if (isModified) modifiedAtAll = true;
      const { type, name, preview, standalone, onlyNew } = fields[field];
      if (onlyNew && !isNew) return; // continue
      else if (standalone) {
        form.push([left, right]);
        left = [];
        right = [];
      }

      if (type === 'date') {
        left.push(
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
        left.push(
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
        left.push(
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
        left.push(
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
        left.push(
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

      if (preview === 'image') {
        if (values[field]) {
          right.push(
            <Img
              key={field}
              className={this.cn({ d: 'rightImg' })}
              src={values[field]}
            />,
          );
        }
      } else if (preview === 'title') {
        right.push(<Markdown key={field} source={`# ${values[field]}`} />);
      } else if (preview === 'text') {
        right.push(<Markdown key={field} source={values[field]} />);
      }

      if (standalone) {
        form.push([left, right]);
        left = [];
        right = [];
      }
    });

    const submitButton = (
      <Button
        key="submit-button"
        type="important"
        text={callForAction || (isNew ?
          `create new ${resourceName || 'resource'}` :
          'save changes'
        )}
        onClick={this.submit}
        disabled={!isNew && !modifiedAtAll}
      />
    );

    form.push([left, right]);
    if (left.length || form.length === 1) left.push(submitButton);
    else form[form.length - 2][0].push(submitButton);

    return (
      <div className={this.cni()}>
        <div className={this.cn({ d: 'inner' })}>
          {form.map(([l, r], i) => (
            // eslint-disable-next-line react/no-array-index-key
            <div key={i} className={this.cn({ d: 'section' })}>
              <div className={this.cn({ d: 'left' })}>{l}</div>
              <div className={this.cn({ d: 'right' })}>{r}</div>
            </div>
          ))}
        </div>
      </div>
    );
  }
}
