import { PropTypes } from 'react';

import Component from '../utils/Component';

export default class Textarea extends Component {
  static propTypes = {
    label: PropTypes.string,
    value: PropTypes.string,
    onChange: PropTypes.func,
    onEnter: PropTypes.func,
    parentField: PropTypes.string, // stateless
    modified: PropTypes.bool,
  };

  constructor(props) {
    super(props);
    this.state = { value: props.value || '' };
  }

  handleChange = (e) => {
    this.setState({ value: e.target.value });
    this.props.onChange && this.props.onChange(
      e.target.value, this.props.parentField,
    );
  };

  handleKeyUp = e =>
    this.props.onEnter && e.which === 13 && this.props.onEnter();

  focus() {
    this.r.input.focus();
  }

  reset() {
    this.setState({ value: '' });
  }

  render() {
    const stateValue = this.state.value;
    return (
      <div className={this.cni()}>
        <textarea
          ref={this.rcb('input')}
          className={this.cn({
            d: 'input',
            s: [!stateValue.length && 'empty'],
          })}
          placeholder={this.props.label}
          value={this.props.parentField ? this.props.value : stateValue}
          onChange={this.handleChange}
          onKeyUp={this.handleKeyUp}
        />
        <div className={this.cn({
          d: 'label',
          s: [!stateValue.length && 'collapsed'],
        })}>
          {this.props.label}
        </div>
        <div
          className={this.cn({
            d: 'modified',
            s: [!this.props.modified && 'collapsed'],
          })}
          title="modified"
        />
        <div className={this.cn({
          d: 'resizer',
          s: [!stateValue.length && 'empty'],
        })}>
          {stateValue}{stateValue[stateValue.length - 1] === '\n' && ' '}
        </div>
      </div>
    );
  }
}
