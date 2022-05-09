public class BankAccount {
    private double balance;
    private double minBalance;
    private boolean isActive;

    public BankAccount(double balance, double minBalance) {
        this.balance = balance;
        this.minBalance = minBalance;
        this.isActive = true;
    }

    public boolean isActive() {
        return isActive;
    }

    public double getBalance() {
        return balance;
    }

    public double getMinBalance() {
        return minBalance;
    }

    public double deposit(double amount) {
        return balance += amount;
    }

    public double withdraw(double amount) {
        if (balance - amount >= minBalance) {
            return balance -= amount;
        }
        else {
            throw new RuntimeException();
        }
    }

}
