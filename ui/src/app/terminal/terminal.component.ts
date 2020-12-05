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

  currentValue: string;
  history: string[] = [];
  historyCursor = -1;
  fitAddon: FitAddon;
  constructor(private api: ApiService) { }

  ngAfterViewInit(): void {
    this.fitAddon = new FitAddon();
    this.terminal.underlying.focus();
    this.terminal.underlying.loadAddon(new WebLinksAddon());
    this.terminal.underlying.loadAddon(this.fitAddon);
    this.terminal.underlying.setOption('cursorBlink', true);
    this.terminal.underlying.setOption('convertEol', true);
    this.terminal.underlying.setOption('theme', {
      background: '#000000',
      cursor: '#99FFFF',
      foreground: '#99FFFF',
    });
    this.api.executeCommand('bio --help', this.terminal.underlying.cols).subscribe(resp => {
      this.terminal.write(`\r\n${resp.output}\n`);
      this.fitAddon.fit();
      this.initTerminal();
    }, err => {
      this.initTerminal();
      console.log(err);
    });
  }

  initTerminal(): void {
    this.currentValue = '';
    this.terminal.write('$ ');

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
      } else if (ev.key === 'ArrowLeft' || ev.key === 'ArrowRight') {
        ev.preventDefault();
      } else if (ev.key === 'ArrowDown' ) {
        ev.preventDefault();
      } else if (ev.key === 'ArrowUp') {
        ev.preventDefault();
        if (this.historyCursor >= 0) {
          //this.currentValue;
        }
      } else if (printable) {
        this.terminal.write(e.key);
        this.currentValue += `${e.key}`;
      }
    });
  }

  lineBreak(clear = false, prefix = true): void {
    const linePref = prefix ? '$ ' : '';
    if (clear) {
      this.currentValue = '';
    }
    this.terminal.write(`\r\n${linePref}`);
  }

  validateCommand(): boolean {
    return this.currentValue.startsWith('clear') || this.currentValue.startsWith('bio ') || this.currentValue === 'bio';
  }

  execute(): void {
    this.currentValue = this.currentValue.trim();
    this.history.push(this.currentValue);
    this.historyCursor = -1;
    if (this.currentValue === 'clear') {
      this.terminal.underlying.reset();
      this.terminal.write('$ ');
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
