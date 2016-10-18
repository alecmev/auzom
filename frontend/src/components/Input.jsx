import { PropTypes } from 'react';

import Component from '../utils/Component';

export default class Input extends Component {
  static propTypes = {
    type: PropTypes.oneOf(['text', 'email', 'password']),
    label: PropTypes.string,
    value: PropTypes.string,
    onChange: PropTypes.func,
    onEnter: PropTypes.func,
    parentField: PropTypes.string, // stateless
    modified: PropTypes.bool,
  };

  static defaultProps = {
    type: 'text',
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

  render() {
    const stateValue = this.state.value;
    return (
      <div className={this.cni()}>
        <input
          ref={this.rcb('input')}
          className={this.cn({
            d: 'input',
            s: [!stateValue.length && 'empty'],
          })}
          type={this.props.type}
          placeholder={this.props.label}
          value={this.props.parentField ? this.props.value : stateValue}
          spellCheck="false"
          onChange={this.handleChange}
          onKeyUp={this.handleKeyUp}
        />
        <div className={this.cn({
          d: 'label',
          s: [!stateValue.length && 'collapsed'],
        })}>
          {this.props.label}
        </div>
        <div className={this.cn({
          d: 'modified',
          s: [!this.props.modified && 'collapsed'],
        })} title="modified" />
      </div>
    );
  }
}
