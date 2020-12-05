export class Help {
  commands: Command[];
}

export class Command {
  command: string;
  description: string;
  aliases: string[];
  flags: Parameters[];
  actions: string;
}

export class Parameters {
  parameter: string;
  type: ParameterType;
  validValues: [];
}

export enum ParameterType {
  STRING = 'string',
  LIST = 'list',
  NUMBER = 'number',
  BOOLEAN = 'boolean',
}
