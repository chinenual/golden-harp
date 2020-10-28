int showVersion(int argc=0, char**argv=NULL)
{
  Serial.println(F( "Harp Controler Version " __DATE__ " " __TIME__ ));
};

int getConfig(int argc, char **argv) {
  config_print();
}

void shell_setup() {
  shell.attach(Serial);
  shell.addCommand(F("version"), showVersion);
  shell.addCommand(F("getconfig"), getConfig);
  showVersion();
}
