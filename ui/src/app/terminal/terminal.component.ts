import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { NgTerminal } from 'ng-terminal';
import { ApiService } from '../api.service';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';

@Component({
  selector: 'app-terminal',
  templateUrl: './terminal.component.html',
  styleUrls: ['./terminal.component.scss']
})
export class TerminalComponent implements AfterViewInit {
  @ViewChild('term', { static: true }) terminal: NgTerminal;

  currentValue = '';
  history: string[] = [];
  maxHistorySize = 50;
  historyCursor = -1;
  fitAddon: FitAddon;
  autocompleteVocabulary = [
    'bio', 'about', 'about-tool', 'certifications', 'education', 'experience', 'help', 'open', 'picture', 'skills',
    '--pretty', '--roles', '--help',
  ];
  ignore = ['ArrowLeft', 'ArrowRight', 'Home', 'End'];
  constructor(private api: ApiService) { }

  ngAfterViewInit(): void {
    this.fitAddon = new FitAddon();
    this.terminal.underlying.focus();
    this.terminal.underlying.loadAddon(new WebLinksAddon());
    this.terminal.underlying.loadAddon(this.fitAddon);
    this.terminal.underlying.setOption('cursorBlink', true);
    this.terminal.underlying.setOption('convertEol', true);
    this.terminal.underlying.setOption('fontWeight', '200');
    this.terminal.underlying.setOption('cursorStyle', 'block');
    this.terminal.underlying.setOption('theme', {
      background: '#111',
      cursor: '#00b3b3',
      foreground: '#00b3b3',
    });
    this.api.executeCommand('bio --help', this.terminal.underlying.cols).subscribe(resp => {
      this.terminal.write(`\r\n${resp.output}\n`);
      this.fitAddon.fit();
      this.initTerminal();
    }, err => {
      this.initTerminal();
    });
  }

  initTerminal(): void {
    this.history = JSON.parse(localStorage.getItem('historyCache') || '[]');
    this.currentValue = '';
    this.terminal.write('➜ ');

    this.terminal.keyEventInput.subscribe(e => {

      const ev = e.domEvent;
      const printable = !ev.altKey && !ev.ctrlKey && !ev.metaKey;

      if (ev.key === 'Enter') {
        this.execute();
      } else if (ev.key === 'Backspace') {
        this.currentValue = this.currentValue.slice(0, -1);
        // Do not delete the prompt
        if (this.terminal.underlying.buffer.active.cursorX > 2) {
          this.terminal.write('\b \b');
        }
      } else if (ev.key === 'Tab') {
        ev.preventDefault();
        const parts = this.currentValue.split(' ');
        const prefix = parts[parts.length - 1];
        if (this.currentValue.trim()) {
          for (const word of this.autocompleteVocabulary) {
            if (word.startsWith(prefix)) {
              const suffix = word.slice(prefix.length, word.length);
              this.currentValue = `${this.currentValue}${suffix}`;
              this.terminal.write(suffix);
              break;
            }
          }
        }


      } else if (this.ignore.indexOf(ev.key) >= 0) {
        ev.preventDefault();
      } else if (ev.key === 'ArrowDown' ) {
        ev.preventDefault();
        this.clearCurrentPrompt();
        this.historyCursor = Math.max(this.historyCursor - 1, -1);
        if (this.historyCursor >= 0) {
          const v = this.history[this.historyCursor];
          this.currentValue = v;
          this.terminal.write(v);
        }
      } else if (ev.key === 'ArrowUp') {
        ev.preventDefault();
        this.clearCurrentPrompt();
        this.historyCursor = Math.min(this.historyCursor + 1, this.history.length - 1);
        const v = this.history[this.historyCursor];
        this.currentValue = v;
        this.terminal.write(v);
      } else if (printable) {
        this.terminal.write(e.key);
        this.currentValue += `${e.key}`;
      }
    });
  }

  clearCurrentPrompt(): void {
    this.terminal.write('\b \b'.repeat(this.currentValue.length));
    this.currentValue = '';
  }

  lineBreak(clear = false, prefix = true): void {
    const linePref = prefix ? '➜ ' : '';
    if (clear) {
      this.currentValue = '';
    }
    this.terminal.write(`\r\n${linePref}`);
  }

  execute(): void {
    this.currentValue = this.currentValue.trim();
    if (this.currentValue !== '') {
      this.history.unshift(this.currentValue);
      this.history = this.history.slice(0, this.maxHistorySize);
      localStorage.setItem('historyCache', JSON.stringify(this.history));
    }
    this.historyCursor = -1;
    if (this.currentValue === 'clear') {
      this.terminal.underlying.reset();
      this.terminal.write('➜ ');
      this.currentValue = '';
    } else if (!this.currentValue.trim()) {
      this.lineBreak(true);
    } else {
      this.lineBreak(false, false);
      this.api.executeCommand(this.currentValue, this.terminal.underlying.cols).subscribe(resp => {
        this.terminal.write(`\r\n${resp.output}`);
        this.fitAddon.fit();
        this.lineBreak(true);
        if (resp.output.startsWith('Opening ')) {
          window.open(resp.output.split(' ')[1], '_blank');
        }
      }, err => {
        const errMessage = typeof err === 'string' ? err : err.message;
        console.log(err);

        this.currentValue = '';
        this.terminal.write('Error: ' + errMessage);
        this.lineBreak(true);
      });
    }
    return;
  }

  clear(): void {
    this.terminal.setStyle(1);
  }
}
